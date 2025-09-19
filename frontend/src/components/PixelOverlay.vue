<template>
  <div 
    class="pixel-overlay" 
    v-if="visible"
    :style="{
      left: `${overlayX}px`,
      top: `${overlayY}px`
    }"
  >
    <div class="window-frame window-outer"></div>
    <div class="window-frame window-inner"></div>
    <div class="center-dot"></div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'PixelOverlay',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    overlayX: {
      type: Number,
      default: 0
    },
    overlayY: {
      type: Number,
      default: 0
    }
  }
});
</script>

<style scoped>
.pixel-overlay {
  position: absolute;
  pointer-events: none;
  z-index: 100;
  transform: translate(-50%, -50%);
}

.window-frame {
  position: absolute;
  border: 2px solid;
  background: transparent;
}

.window-outer {
  width: 13px;
  height: 13px;
  border-color: black;
  transform: translate(-50%, -50%); /* Центрирование рамки */
}

.window-inner {
  width: 9px;
  height: 9px;
  border-color: white;
  transform: translate(-50%, -50%); /* Центрирование рамки */
}
</style>