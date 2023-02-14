<template>
    <div class="container">
        <div>
            <h3 class="close-button" @click="close()"> close &times; </h3>            
        </div>
        <slot/>
    </div>                      
</template>

<script>
export default {
    methods: {
        close() {
            this.$emit("close")
        }
    }

}
</script>

<style>
@media only screen and (min-width: 1000px) {
    .container {
        width: 60rem;
        margin: auto
    }
}

@media only screen and (max-width: 1000px) {
    .container {
        
        width: auto;
        margin-right: 1rem;
        margin-left: 1rem;
    }
}

.container {
    background-color: white;
    border-radius: 8px;
    height: 40rem;
    margin-top: 10rem;
    margin-bottom: auto;
    
}
h3.close-button {
    float: right;
    font-size: 1rem;
    background: none;
    border: none;
    margin-right: 0.5rem;
    margin-top: 0.5rem;
    cursor: pointer;
    padding: 0.3rem;
    height: 0.8rem;
}

h3.close-button:hover {
    background: rgb(199, 199, 199);
    border-radius: 9999px;
}

</style>