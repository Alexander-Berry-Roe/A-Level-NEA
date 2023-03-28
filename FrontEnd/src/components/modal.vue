<template>
  <div class="modal-backdrop">
    <div class="modal">
      <slot name="body">
      </slot>-
    </div>
  </div>
</template>
<script>
  export default {
    name: 'Modal',
    methods: {
      close() {
        this.$emit('close');
      },
    },
  };
</script>
<style scoped>
.modal-backdrop {
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    align-items: center;
    background-color: rgba(0, 0, 0, 0.3);
    display: flex;
    justify-content: center;
  }

.modal {
  position: fixed;
  z-index: 100;
  background: #FFFFFF;
  box-shadow: 2px 2px 20px 1px;
  overflow-x: auto;
  display: flex;
  flex-direction: column;
  pointer-events: auto; 
}


</style>