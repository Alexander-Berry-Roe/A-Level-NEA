<template>
    <div class="message">
        <h1>
            Are you sure you want to log out?
        </h1>
        <button class="button-yes" @click="sendLogout()"> Yes </button> 
        <button class="button-no" @click="close()"> No </button>
    </div>
</template>

<script>
import axios from "axios";
export default {
    data () {
        return {}
    },
    methods: {
        close() {
            this.$emit("close")
        },
        sendLogout() {
            axios
                .get("/api/logout")
                .then(response => {
                    window.location.reload()
                })
        }

    }
}
</script>

<style>
.message {
    background-color: white;
    width: 15rem;
    font-family: 'Courier New', Courier, monospace;
    font-size: 0.5rem;
    margin: auto;
    margin-top: 20rem;
    border-radius: 5px;
    padding: 1rem;
    z-index: 100;
}
.button-yes {
    font-family: 'Courier New', Courier, monospace;
    color: white;
    margin: auto;
    outline: none;
    width: 6.5rem;
    height: 2rem;
    background: rgb(255, 0, 0);
    border-radius: 5px;
    border: none;
    cursor: pointer;
}
.button-yes:hover {
    background: rgb(126, 126, 126);
}
.button-no {
    font-family: 'Courier New', Courier, monospace;
    color: white;
    margin: auto;
    outline: none;
    width: 6.5rem;
    height: 2rem;
    float: right;
    background: rgb(0, 102, 255);
    border-radius: 5px;
    border: none;
    cursor: pointer;
}
.button-no:hover{
    background: rgb(0, 92, 225);
}
</style>