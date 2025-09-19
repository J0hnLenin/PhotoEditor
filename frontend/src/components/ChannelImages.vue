<template>
  <div class="channel-images">
    <h3>Каналы изображения</h3>
    <div class="channels-grid">
      <div 
        v-for="channel in channels" 
        :key="channel.type"
        class="channel-item"
      >
        <img 
          :src="channel.url" 
          :alt="`${channel.type} канал`"
          class="channel-image"
        />
        <span class="channel-label">{{ getChannelLabel(channel.type) }}</span>
      </div>
    </div>
    <div v-if="loading" class="channels-loading">Загрузка каналов...</div>
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue';
import type { ChannelImage } from '@/types/image';

export default defineComponent({
  name: 'ChannelImages',
  props: {
    channels: {
      type: Array as PropType<ChannelImage[]>,
      required: true
    },
    loading: {
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