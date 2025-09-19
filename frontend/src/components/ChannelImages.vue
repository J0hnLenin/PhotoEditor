<template>
  <div class="channel-images">
    <h3>Каналы изображения </h3>
    <div class="channels-grid">
      <div 
        v-for="channel in channels" 
        :key="channel.type"
        class="channel-item"
      >
        <div class="channel-image-container">
          <img 
            :src="channel.url" 
            :alt="`${channel.type} канал`"
            class="channel-image"
          />
          <span class="channel-label">{{ getChannelLabel(channel.type) }}</span>
        </div>
        
        <Histogram 
          v-if="channel.histogram"
          :data="channel.histogram"
          :color="getChannelColor(channel.type)"
          :width="150"
          :height="100"
        />
        
        <div v-else class="histogram-placeholder">
          Загрузка гистограммы...
        </div>
      </div>
    </div>
    <div v-if="loading" class="channels-loading">Загрузка каналов...</div>
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue';
import Histogram from '@/components/Histogram.vue';
import type { ChannelImage } from '@/types/image';

export default defineComponent({
  name: 'ChannelImages',
  components: {
    Histogram
  },
  props: {
    channels: {
      type: Array as PropType<ChannelImage[]>,
      required: true
    },
    loading: {
      type: Boolean,
      default: false
    },
    processedImage: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    getChannelLabel(type: ChannelImage['type']): string {
      const labels = {
        red: 'Красный',
        green: 'Зеленый',
        blue: 'Синий',
        gray: 'Серый'
      };
      return labels[type];
    },
    
    getChannelColor(type: ChannelImage['type']): string {
      const colors = {
        red: '#dc3545',
        green: '#28a745',
        blue: '#007bff',
        gray: '#6c757d'
      };
      return colors[type];
    }
  }
});
</script>

<style scoped>
.channel-images {
  border-top: 1px solid #eee;
  padding-top: 20px;
}

.channel-images h3 {
  margin-bottom: 15px;
  color: #333;
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.image-type {
  font-size: 12px;
  color: #666;
  font-weight: normal;
  font-style: italic;
}

.channels-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
}

.channel-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.channel-image-container {
  position: relative;
  border: 1px solid #ddd;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 8px;
}

.channel-image {
  width: 100%;
  height: 80px;
  object-fit: cover;
  display: block;
}

.channel-label {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 4px;
  font-size: 12px;
  text-align: center;
}

.histogram-placeholder {
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  color: #999;
  background: #f8f9fa;
  border-radius: 4px;
  width: 100%;
}

.channels-loading {
  text-align: center;
  padding: 20px;
  color: #666;
  font-style: italic;
}

@media (max-width: 768px) {
  .channels-grid {
    grid-template-columns: 1fr;
  }
}

.channel-images {
  border-top: 1px solid #eee;
  padding-top: 20px;
}

.channel-images h3 {
  margin-bottom: 15px;
  color: #333;
  font-size: 16px;
}

.channels-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.channel-item {
  position: relative;
  border: 1px solid #ddd;
  border-radius: 4px;
  overflow: hidden;
}

.channel-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.channel-label {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 4px;
  font-size: 12px;
  text-align: center;
}

.channels-loading {
  text-align: center;
  padding: 20px;
  color: #666;
  font-style: italic;
}
</style>