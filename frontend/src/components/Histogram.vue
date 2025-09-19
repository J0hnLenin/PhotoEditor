<template>
  <div class="histogram">
    <div class="histogram-title">Гистограмма яркости</div>
    <div class="histogram-container" ref="histogramContainer">
      <svg :width="width" :height="height" class="histogram-svg">
        <g v-for="(value, index) in data" :key="index">
          <rect 
            :x="index * barWidth" 
            :y="height - (value * scaleFactor)" 
            :width="barWidth - 0.5" 
            :height="value * scaleFactor" 
            :fill="color"
            class="histogram-bar"
            :opacity="0.8"
          />
        </g>
      </svg>
    </div>
    
    <div class="histogram-labels">
      <h2>0</h2>
      <h2>255</h2>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType, ref, computed, onMounted, watch } from 'vue';

export default defineComponent({
  name: 'Histogram',
  props: {
    data: {
      type: Array as PropType<number[]>,
      required: true
    },
    color: {
      type: String,
      default: '#007bff'
    },
    width: {
      type: Number,
      default: 1000
    },
    height: {
      type: Number,
      default: 500
    }
  },
  setup(props) {
    const histogramContainer = ref<HTMLElement | null>(null);
    const containerWidth = ref(props.width);

    const barWidth = computed(() => {
      return containerWidth.value / 256;
    });

    const maxValue = computed(() => {
      return Math.max(...props.data);
    });

    const scaleFactor = computed(() => {
      return maxValue.value > 0 ? props.height / maxValue.value : 0;
    });

    onMounted(() => {
      if (histogramContainer.value) {
        containerWidth.value = histogramContainer.value.clientWidth;
      }
    });

    watch(() => props.data, () => {
      if (histogramContainer.value) {
        containerWidth.value = histogramContainer.value.clientWidth;
      }
    });

    return {
      histogramContainer,
      barWidth,
      scaleFactor
    };
  }
});
</script>

<style scoped>
.histogram {
  margin-top: 12px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  width: 100%;
}

.histogram-title {
  font-size: 12px;
  font-weight: 600;
  color: #495057;
  margin-bottom: 8px;
  text-align: center;
}

.histogram-container {
  width: 100%;
  height: v-bind('height + "px"');
  overflow: hidden;
  background: #ffffff;
  border-radius: 4px;
  border: 1px solid #dee2e6;
  padding: 4px;
}

.histogram-svg {
  display: block;
  width: 100%;
  height: 100%;
}

.histogram-bar {
  transition: all 0.2s ease;
}

.histogram-bar:hover {
  opacity: 1 !important;
  filter: brightness(1.2);
}

.histogram-labels {
  display: flex;
  justify-content: space-between;
  margin-top: 6px;
  font-size: 10px;
  color: #6c757d;
}

</style>