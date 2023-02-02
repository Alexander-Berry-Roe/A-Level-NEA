<template>
    <div class="modal" >
        <slot/>
    </div>
</template>

<script>
export default {
    data() {
        return {

        }
    },
    methods: {
        close() {
            this.$emit("close")
        },
        log() {
            console.log("click")
        }
    }
}
</script>

<style>
.modal {
  overflow-x: hidden;
  overflow-y: auto;
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 9;
    background-color: rgba(0, 0, 0, 0.3);
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    z-index: 1;
}
</style>