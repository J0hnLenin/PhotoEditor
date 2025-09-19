<template>
  <div class="main-image">
    <div 
      class="image-container"
      ref="imageContainer"
      @mousemove="handleMouseMove"
      @mouseleave="handleMouseLeave"
    >
      <img 
        :src="image" 
        alt="Основное изображение"
        class="image"
        :class="{ loading: loading }"
        ref="imageElement"
        @load="onImageLoad"
      />
      
      <PixelOverlay 
        :visible="showOverlay"
        :overlay-x="overlayX"
        :overlay-y="overlayY"
      />
      
      <div v-if="loading" class="loader">Обработка...</div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType, ref, onMounted, onUnmounted } from 'vue';
import PixelOverlay from './PixelOverlay.vue';

interface RGB {
  r: number;
  g: number;
  b: number;
}

interface WindowStats {
  mean: number;
  stdDev: number;
}

export default defineComponent({
  name: 'MainImage',
  components: {
    PixelOverlay
  },
  props: {
    image: {
      type: String as PropType<string | null>,
      required: true
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  emits: ['pixel-data'],
  setup(props, { emit }) {
    const imageContainer = ref<HTMLElement | null>(null);
    const imageElement = ref<HTMLImageElement | null>(null);
    const showOverlay = ref(false);
    const overlayX = ref(0);
    const overlayY = ref(0);
    
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');

    const onImageLoad = () => {
      if (imageElement.value && ctx) {
        canvas.width = imageElement.value.naturalWidth;
        canvas.height = imageElement.value.naturalHeight;
        ctx.drawImage(imageElement.value, 0, 0);
      }
    };

    const getPixelData = (x: number, y: number): RGB => {
      if (!ctx) return { r: 0, g: 0, b: 0 };
      
      const imageData = ctx.getImageData(x, y, 1, 1).data;
      return {
        r: imageData[0],
        g: imageData[1],
        b: imageData[2]
      };
    };

    const getWindowData = (centerX: number, centerY: number): RGB[] => {
      if (!ctx) return [];
      
      const windowSize = 11;
      const halfWindow = Math.floor(windowSize / 2);
      const pixels: RGB[] = [];

      for (let y = centerY - halfWindow; y <= centerY + halfWindow; y++) {
        for (let x = centerX - halfWindow; x <= centerX + halfWindow; x++) {
          if (x >= 0 && x < canvas.width && y >= 0 && y < canvas.height) {
            const imageData = ctx.getImageData(x, y, 1, 1).data;
            pixels.push({
              r: imageData[0],
              g: imageData[1],
              b: imageData[2]
            });
          } else {
            // Для пикселей за границами изображения
            pixels.push({ r: 0, g: 0, b: 0 });
          }
        }
      }

      return pixels;
    };

    const calculateWindowStats = (pixels: RGB[]): WindowStats => {
      const intensities = pixels.map(p => (p.r + p.g + p.b) / 3);
      const mean = intensities.reduce((sum, val) => sum + val, 0) / intensities.length;
      
      const variance = intensities.reduce((sum, val) => {
        return sum + Math.pow(val - mean, 2);
      }, 0) / intensities.length;
      
      const stdDev = Math.sqrt(variance);
      
      return { mean, stdDev };
    };

    const handleMouseMove = (event: MouseEvent) => {
      if (!imageElement.value || !imageContainer.value) return;

      const rect = imageElement.value.getBoundingClientRect();
      const containerRect = imageContainer.value.getBoundingClientRect();
      
      // Координаты относительно изображения
      const x = Math.floor((event.clientX - rect.left) * (imageElement.value.naturalWidth / rect.width));
      const y = Math.floor((event.clientY - rect.top) * (imageElement.value.naturalHeight / rect.height));

      // Координаты для overlay
      overlayX.value = event.clientX - containerRect.left;
      overlayY.value = event.clientY - containerRect.top;
      showOverlay.value = true;

      if (x >= 0 && x < imageElement.value.naturalWidth && 
          y >= 0 && y < imageElement.value.naturalHeight) {
        
        const pixelData = getPixelData(x, y);
        const intensity = (pixelData.r + pixelData.g + pixelData.b) / 3;
        
        const windowData = getWindowData(x, y);
        const windowStats = calculateWindowStats(windowData);

        emit('pixel-data', {
          x,
          y,
          rgb: pixelData,
          intensity,
          windowStats,
          windowData
        });
      }
    };

    const handleMouseLeave = () => {
      showOverlay.value = false;
      emit('pixel-data', null);
    };

    onMounted(() => {
      if (imageElement.value?.complete) {
        onImageLoad();
      }
    });

    return {
      imageContainer,
      imageElement,
      showOverlay,
      overlayX,
      overlayY,
      onImageLoad,
      handleMouseMove,
      handleMouseLeave
    };
  }
});
</script>

<style scoped>
.main-image {
  width: 100%;
}

.image-container {
  position: relative;
  border: 2px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
  background: #f8f9fa;
  cursor: none;
}

.image {
  width: 100%;
  height: auto;
  display: block;
  transition: opacity 0.3s ease;
}

.image.loading {
  opacity: 0.5;
}

.loader {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 10px 20px;
  border-radius: 20px;
  font-size: 14px;
}
</style>