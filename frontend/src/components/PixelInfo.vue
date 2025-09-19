<template>
  <div class="pixel-info" v-if="visible">
    <div class="info-section">
      <h4>Информация о пикселе</h4>
      <div class="info-grid">
        <div class="info-item">
          <span class="label">Координаты:</span>
          <span class="value">({{ pixelX }}, {{ pixelY }})</span>
        </div>
        <div class="info-item">
          <span class="label">R:</span>
          <span class="value red">{{ rgb.r }}</span>
        </div>
        <div class="info-item">
          <span class="label">G:</span>
          <span class="value green">{{ rgb.g }}</span>
        </div>
        <div class="info-item">
          <span class="label">B:</span>
          <span class="value blue">{{ rgb.b }}</span>
        </div>
        <div class="info-item">
          <span class="label">Интенсивность:</span>
          <span class="value">{{ intensity.toFixed(1) }}</span>
        </div>
        <div class="info-item">
          <span class="label">Среднее (μ):</span>
          <span class="value">{{ windowStats.mean.toFixed(1) }}</span>
        </div>
        <div class="info-item">
          <span class="label">Станд. отклонение (σ):</span>
          <span class="value">{{ windowStats.stdDev.toFixed(1) }}</span>
        </div>
      </div>
    </div>

    <div class="window-preview" v-if="windowData">
      <h4>Окно 11×11 (увеличено)</h4>
      <div class="window-container">
        <div 
          class="pixel-grid"
          :style="{
            gridTemplateColumns: `repeat(11, 8px)`,
            gridTemplateRows: `repeat(11, 8px)`
          }"
        >
          <div
            v-for="(pixel, index) in windowData"
            :key="index"
            class="pixel"
            :style="{
              backgroundColor: `rgb(${pixel.r}, ${pixel.g}, ${pixel.b})`,
              border: isCenterPixel(index) ? '1px solid #ff0' : '1px solid #ccc'
            }"
            :title="`RGB: ${pixel.r},${pixel.g},${pixel.b}`"
          ></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue';

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
  name: 'PixelInfo',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    pixelX: {
      type: Number,
      default: 0
    },
    pixelY: {
      type: Number,
      default: 0
    },
    rgb: {
      type: Object as PropType<RGB>,
      default: () => ({ r: 0, g: 0, b: 0 })
    },
    intensity: {
      type: Number,
      default: 0
    },
    windowStats: {
      type: Object as PropType<WindowStats>,
      default: () => ({ mean: 0, stdDev: 0 })
    },
    windowData: {
      type: Array as PropType<RGB[]>,
      default: () => []
    }
  },
  methods: {
    isCenterPixel(index: number): boolean {
      // Центральный пиксель в окне 11x11 (индекс 60)
      return index === 60;
    }
  }
});
</script>

<style scoped>
.pixel-info {
  position: fixed;
  top: 20px;
  right: 20px;
  background: white;
  border: 2px solid #007bff;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  min-width: 300px;
  max-width: 400px;
}

.info-section {
  margin-bottom: 20px;
}

.info-section h4 {
  margin: 0 0 12px 0;
  color: #2c3e50;
  font-size: 16px;
  font-weight: 600;
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 0;
  border-bottom: 1px solid #f1f3f4;
}

.label {
  font-weight: 500;
  color: #495057;
  font-size: 12px;
}

.value {
  font-weight: 600;
  color: #2c3e50;
  font-size: 12px;
}

.value.red { color: #dc3545; }
.value.green { color: #28a745; }
.value.blue { color: #007bff; }

.window-preview h4 {
  margin: 0 0 12px 0;
  color: #2c3e50;
  font-size: 16px;
  font-weight: 600;
}

.window-container {
  background: #f8f9fa;
  padding: 12px;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  display: flex;
  justify-content: center;
}

.pixel-grid {
  display: grid;
  gap: 1px;
  background: #ccc;
  padding: 1px;
  border: 1px solid #999;
}

.pixel {
  width: 8px;
  height: 8px;
  cursor: pointer;
  transition: transform 0.1s ease;
}

.pixel:hover {
  transform: scale(1.5);
  z-index: 2;
  box-shadow: 0 0 4px rgba(0, 0, 0, 0.5);
}

@media (max-width: 768px) {
  .pixel-info {
    top: 10px;
    right: 10px;
    left: 10px;
    max-width: none;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .pixel-grid {
    transform: scale(0.8);
  }
}
</style>