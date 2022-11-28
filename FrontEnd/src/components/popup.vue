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
        }
    }
}
</script>

<style>
.modal {
  position: absolute;
  margin: 0px;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 99;
  background: transparent;

  
}
</style>