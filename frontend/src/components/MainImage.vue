<template>
  <div class="main-image">
    <div class="image-container">
      <img 
        :src="image" 
        alt="Основное изображение"
        class="image"
        :class="{ loading: loading }"
      />
      <div v-if="loading" class="loader">Обработка...</div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue';

export default defineComponent({
  name: 'MainImage',
  props: {
    image: {
      type: String as PropType<string | null>,
      required: true
    },
    loading: {
      type: Boolean,
      default: false
    }
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