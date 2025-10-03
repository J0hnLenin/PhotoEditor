<template>
  <div class="controls-panel">
    <div class="panel-header">
      <h3>Настройки изображения</h3>
      <button 
        @click="resetSettings" 
        class="reset-btn"
        :disabled="applying"
        title="Сбросить все настройки"
      >
        ↺ Сбросить
      </button>
    </div>
    
    <div class="control-group">
      <label>Яркость красного: {{ params.RedBrightness }}%</label>
      <input 
        type="range" 
        min="1" 
        max="200" 
        v-model.number="localParams.RedBrightness"
        @input="updateParam('RedBrightness', $event)"
      />
    </div>

    <div class="control-group">
      <label>Яркость зеленого: {{ params.GreenBrightness }}%</label>
      <input 
        type="range" 
        min="1" 
        max="200" 
        v-model.number="localParams.GreenBrightness"
        @input="updateParam('GreenBrightness', $event)"
      />
    </div>

    <div class="control-group">
      <label>Яркость синего: {{ params.BlueBrightness }}%</label>
      <input 
        type="range" 
        min="1" 
        max="200" 
        v-model.number="localParams.BlueBrightness"
        @input="updateParam('BlueBrightness', $event)"
      />
    </div>

    <div class="control-group">
      <label>Контраст: {{ params.Contrast }} %</label>
      <input 
        type="range" 
        min="1" 
        max="200" 
        v-model.number="localParams.Contrast"
        @input="updateParam('Contrast', $event)"
      />
    </div>

    <div class="control-group">
      <label>Magic: {{ params.Magic }} %</label>
      <input 
        type="range" 
        min="0" 
        max="150" 
        v-model.number="localParams.Magic"
        @input="updateParam('Magic', $event)"
      />
    </div>

    <div class="control-group checkbox">
      <label>
        <input 
          type="checkbox" 
          v-model="localParams.Negative"
          @change="updateParam('Negative', $event)"
        />
        Негатив
      </label>
    </div>

    <div class="control-group checkbox">
      <label>
        <input 
          type="checkbox" 
          v-model="localParams.VerticalMirror"
          @change="updateParam('VerticalMirror', $event)"
        />
        Вертикальное зеркало
      </label>
    </div>

    <div class="control-group checkbox">
      <label>
        <input 
          type="checkbox" 
          v-model="localParams.HorizontalMirror"
          @change="updateParam('HorizontalMirror', $event)"
        />
        Горизонтальное зеркало
      </label>
    </div>

    <div class="control-group">
      <label>Порядок каналов:</label>
      <select 
        v-model="localParams.Order"
        @change="updateParam('Order', $event)"
      >
        <option value="RGB">RGB</option>
        <option value="RBG">RBG</option>
        <option value="GRB">GRB</option>
        <option value="GBR">GBR</option>
        <option value="GBR">BGR</option>
        <option value="GBR">BRG</option>
      </select>
    </div>
    <div class="control-group">
      <label>Фильтр:</label>
      <select 
        v-model="localParams.Filter"
        @change="updateParam('Filter', $event)"
      >
        <option value="none">Без фильтра</option>
        <option value="gaussian">Гауссово размытие</option>
        <option value="median">Медианный фильтр</option>
        <option value="sigma">Сигма фильтр</option>
      </select>
    </div>

    <div v-if="localParams.Filter !== 'none'" class="filter-params">
      <div class="control-group">
        <label>Размер ядра: {{ localParams.FilterSize }}×{{ localParams.FilterSize }}</label>
        <select 
          v-model="localParams.FilterSize"
          @change="updateParam('FilterSize', $event)"
        >
          <option value="3">3×3</option>
          <option value="5">5×5</option>
          <option value="7">7×7</option>
          <option value="9">9×9</option>
          <option value="11">11×11</option>
        </select>
      </div>
    </div>
    </div>
</template>

<script lang="ts">
import { defineComponent, type PropType, reactive, watch } from 'vue';
import type { ImageEditorParams } from '@/types/image';

// Значения по умолчанию для сброса
const DEFAULT_PARAMS: ImageEditorParams = {
  RedBrightness: 100,
  GreenBrightness: 100,
  BlueBrightness: 100,
  Contrast: 100,
  Negative: false,
  Order: 'RGB',
  VerticalMirror: false,
  HorizontalMirror: false,
  Magic: 0,
  Filter: 'none',
  FilterSize: 3,
  Sigma: 0
};

export default defineComponent({
  name: 'ControlsPanel',
  props: {
    params: {
      type: Object as PropType<ImageEditorParams>,
      required: true
    },
    applying: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:params', 'reset'],
  setup(props, { emit }) {
    const localParams = reactive({ ...props.params });

    watch(() => props.params, (newParams) => {
      Object.assign(localParams, newParams);
    });

    const updateParam = (key: keyof ImageEditorParams, event: Event) => {
      const target = event.target as HTMLInputElement | HTMLSelectElement;
      let value: any = target.value;
      
      if (target.type === 'checkbox') {
        value = (target as HTMLInputElement).checked;
      } else if (target.type === 'range') {
        value = Number(value);
      }

      emit('update:params', { [key]: value });
    };

    const resetSettings = () => {
      Object.assign(localParams, DEFAULT_PARAMS);
      emit('reset');
      emit('update:params', DEFAULT_PARAMS);
    };

    return {
      localParams,
      updateParam,
      resetSettings
    };
  }
});
</script>

<style scoped>
.controls-panel {
  background: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  border: 1px solid #dee2e6;
  position: relative;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.panel-header h3 {
  margin: 0;
  color: #333;
  font-size: 18px;
}

.reset-btn {
  padding: 6px 12px;
  background: #6c757d;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  transition: background-color 0.2s;
}

.reset-btn:hover:not(:disabled) {
  background: #5a6268;
}

.reset-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.status {
  position: absolute;
  top: 10px;
  right: 20px;
  font-size: 12px;
  color: #666;
}

.loading-dots::after {
  content: '...';
  animation: dots 1.5s steps(5, end) infinite;
}

@keyframes dots {
  0%, 20% { content: '.'; }
  40% { content: '..'; }
  60%, 100% { content: '...'; }
}

.auto-apply-note {
  margin-bottom: 15px;
  padding: 8px;
  background: #e9ecef;
  border-radius: 4px;
  text-align: center;
}

.auto-apply-note small {
  color: #6c757d;
  font-style: italic;
}

.control-group {
  margin-bottom: 15px;
}

.control-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
  color: #555;
}

.control-group input[type="range"] {
  width: 100%;
  height: 6px;
  border-radius: 3px;
  background: #ddd;
  outline: none;
  -webkit-appearance: none;
}

.control-group input[type="range"]::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #007bff;
  cursor: pointer;
}

.control-group.checkbox {
  display: flex;
  align-items: center;
}

.control-group.checkbox label {
  display: flex;
  align-items: center;
  margin-bottom: 0;
}

.control-group.checkbox input {
  margin-right: 8px;
}

.control-group select {
  width: 100%;
  padding: 8px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  background: white;
}
</style>