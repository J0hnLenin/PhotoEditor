<template>
  <div class="image-editor">
    <div class="header">
      <input 
        type="file" 
        accept=".png" 
        @change="handleFileUpload" 
        ref="fileInput"
        class="file-input"
      />
      
      <div class="action-buttons">
        <button @click="triggerFileInput" class="upload-btn">
        Загрузить PNG изображение
        </button>
        <button v-if="originalImage"
          @click="downloadImage" 
          class="action-btn download-btn"
          :disabled="!processedImage || processing"
        >
          Скачать
        </button>
      </div>
    </div>

    <div class="editor-content" v-if="originalImage">
      <div class="main-section">
        <MainImage 
          :image="processedImage || originalImage"
          :loading="processing"
          @pixel-data="handlePixelData"
        />
        
        <ChannelImages 
          :channels="channelImages"
          :loading="channelsLoading"
          :processed-image="!!processedImage"
        />
      </div>

      <ControlsPanel 
        :params="editorParams"
        @update:params="updateParams"
        @reset="resetSettings"
        :applying="processing"
      />

      <PixelInfo 
        :visible="!!pixelData"
        :pixel-x="pixelData?.x || 0"
        :pixel-y="pixelData?.y || 0"
        :rgb="pixelData?.rgb || { r: 0, g: 0, b: 0 }"
        :intensity="pixelData?.intensity || 0"
        :window-stats="pixelData?.windowStats || { mean: 0, stdDev: 0 }"
        :window-data="pixelData?.windowData || []"
      />
    </div>

    <div v-else class="placeholder">
      <p>Загрузите PNG изображение для начала редактирования</p>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, reactive, watch, onUnmounted } from 'vue';
import MainImage from '@/components/MainImage.vue';
import ChannelImages from '@/components/ChannelImages.vue';
import ControlsPanel from '@/components/ControlsPanel.vue';
import PixelInfo from '@/components/PixelInfo.vue';
import type { ImageEditorParams, ChannelImage, BrightnessHistogram } from '@/types/image';

interface PixelData {
  x: number;
  y: number;
  rgb: { r: number; g: number; b: number };
  intensity: number;
  windowStats: { mean: number; stdDev: number };
  windowData: Array<{ r: number; g: number; b: number }>;
}

// Значения по умолчанию
export const DEFAULT_PARAMS: ImageEditorParams = {
  RedBrightness: 100,
  GreenBrightness: 100,
  BlueBrightness: 100,
  Contrast: 100,
  Negative: false,
  Order: 'RGB',
  VerticalMirror: false,
  HorizontalMirror: false,
  Magic: 0,
  Filter:  'none',
  FilterSize: 3,
  Sigma: 0,
  Interval: 1,
  UnsharpMasking: 0,
};

export default defineComponent({
  name: 'ImageEditor',
  components: {
    MainImage,
    ChannelImages,
    ControlsPanel,
    PixelInfo,
  },
  setup() {
    const fileInput = ref<HTMLInputElement | null>(null);
    const originalImage = ref<string | null>(null);
    const originalFile = ref<File | null>(null);
    const processedImage = ref<string | null>(null);
    const processedImageBlob = ref<Blob | null>(null);
    const channelImages = ref<ChannelImage[]>([]);
    const channelsLoading = ref(false);
    const processing = ref(false);
    const currentStatistics = ref<BrightnessHistogram | null>(null);
    
    let applyTimer: number | null = null;

    const editorParams = reactive<ImageEditorParams>({ ...DEFAULT_PARAMS });

    const pixelData = ref<PixelData | null>(null);

    const handlePixelData = (data: PixelData | null) => {
      pixelData.value = data;
    };

    onUnmounted(() => {
      if (applyTimer) clearTimeout(applyTimer);
    });

    const triggerFileInput = () => {
      fileInput.value?.click();
    };

    const handleFileUpload = (event: Event) => {
      const input = event.target as HTMLInputElement;
      const file = input.files?.[0];
      
      if (file && file.type === 'image/png') {
        originalFile.value = file;
        const reader = new FileReader();
        reader.onload = (e) => {
          originalImage.value = e.target?.result as string;
          processedImage.value = null;
          processedImageBlob.value = null;
          channelImages.value = [];
          currentStatistics.value = null;
          
          // При загрузке нового файла сразу применяем изменения
          scheduleApplyChanges();
        };
        reader.readAsDataURL(file);
      }
    };

    const updateParams = (newParams: Partial<ImageEditorParams>) => {
      Object.assign(editorParams, newParams);
      scheduleApplyChanges();
    };

    const scheduleApplyChanges = () => {
      if (!originalFile.value) return;

      if (applyTimer) clearTimeout(applyTimer);
      
      applyTimer = window.setTimeout(async () => {
        await applyChanges();
      }, 500);
    };

    const applyChanges = async () => {
      if (!originalFile.value) return;

      processing.value = true;
      channelsLoading.value = true;
      try {
        const queryParams = new URLSearchParams();
        Object.entries(editorParams).forEach(([key, value]) => {
          queryParams.append(key, value.toString());
        });

        // Создаем FormData для запроса
        const formData = new FormData();
        formData.append('image', originalFile.value);

        const editorResponse = await fetch(`http://localhost:8000/api/v1/image/redactor?${queryParams}`, {
          method: 'POST',
          body: formData
        });

        if (editorResponse.ok) {
          // Парсим multipart/form-data ответ
          const responseFormData = await editorResponse.formData();
          
          // Получаем все изображения из ответа
          const redactedImageBlob = responseFormData.get('redacted_image') as Blob;
          const redChannelBlob = responseFormData.get('red_channel') as Blob;
          const greenChannelBlob = responseFormData.get('green_channel') as Blob;
          const blueChannelBlob = responseFormData.get('blue_channel') as Blob;
          const grayChannelBlob = responseFormData.get('gray_channel') as Blob;
          const ChangesChannelBlob = responseFormData.get('changes_channel') as Blob;
          
          // Получаем статистику
          const statisticsText = responseFormData.get('statistics') as string;
          const statistics = JSON.parse(statisticsText);

          // Обновляем основное изображение
          processedImageBlob.value = redactedImageBlob;
          processedImage.value = URL.createObjectURL(redactedImageBlob);

          // Обновляем канальные изображения
          channelImages.value = [
            { 
              type: 'red', 
              url: URL.createObjectURL(redChannelBlob), 
              histogram: statistics.Brightness?.Red ? Array.from(statistics.Brightness.Red) : undefined 
            },
            { 
              type: 'green', 
              url: URL.createObjectURL(greenChannelBlob), 
              histogram: statistics.Brightness?.Green ? Array.from(statistics.Brightness.Green) : undefined 
            },
            { 
              type: 'blue', 
              url: URL.createObjectURL(blueChannelBlob), 
              histogram: statistics.Brightness?.Blue ? Array.from(statistics.Brightness.Blue) : undefined 
            },
            { 
              type: 'gray', 
              url: URL.createObjectURL(grayChannelBlob), 
              histogram: statistics.Brightness?.Gray ? Array.from(statistics.Brightness.Gray) : undefined 
            },
            { 
              type: 'changes', 
              url: URL.createObjectURL(ChangesChannelBlob), 
              histogram: undefined  
            },
          ];

          // Сохраняем статистику
          currentStatistics.value = statistics;
        }
      } catch (error) {
        console.error('Ошибка применения изменений:', error);
      } finally {
        processing.value = false;
        channelsLoading.value = false;
      }
    };

    const resetSettings = () => {
      Object.assign(editorParams, DEFAULT_PARAMS);
      processedImage.value = null;
      processedImageBlob.value = null;
      channelImages.value = [];
      currentStatistics.value = null;
    };

    const downloadImage = () => {
      if (!processedImageBlob.value) return;

      const url = URL.createObjectURL(processedImageBlob.value);
      const link = document.createElement('a');
      link.href = url;
      link.download = `edited-image-${new Date().getTime()}.png`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
    };

    watch(editorParams, () => {
      scheduleApplyChanges();
    });

    return {
      fileInput,
      originalImage,
      processedImage,
      channelImages,
      processing,
      channelsLoading,
      editorParams,
      currentStatistics,
      triggerFileInput,
      handleFileUpload,
      updateParams,
      resetSettings,
      downloadImage,
      pixelData,
      handlePixelData
    };
  }
});
</script>

<style scoped>
.image-editor {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  margin-bottom: 30px;
  text-align: center;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 15px;
}

.file-input {
  display: none;
}

.upload-btn {
  padding: 12px 24px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.2s;
}

.upload-btn:hover {
  background: #0056b3;
}

.action-buttons {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

.action-btn {
  padding: 12px 24px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 16px;
  font-weight: 500;
  transition: all 0.2s;
}

.action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.download-btn {
  background: #28a745;
  color: white;
}

.download-btn:hover:not(:disabled) {
  background: #218838;
}

.editor-content {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 30px;
}

.main-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.placeholder {
  text-align: center;
  padding: 100px 0;
  color: #666;
  font-size: 18px;
}

@media (max-width: 768px) {
  .editor-content {
    grid-template-columns: 1fr;
  }
  
  .action-buttons {
    flex-direction: column;
    width: 100%;
    max-width: 300px;
  }
  
  .action-btn {
    width: 100%;
  }
}
</style>