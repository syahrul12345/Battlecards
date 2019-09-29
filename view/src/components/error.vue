<template>
    <v-card class="text-center">
        <v-card-text class="pt-2">{{message}}</v-card-text>
        <v-card-actions class="justify-center">
            <v-btn @click="closeError"> CLOSE </v-btn>
        </v-card-actions>
    </v-card>
</template>

<script>
    export default{
        props:['message'],
        methods: {
            async closeError() {
                this.$emit('closeerror')
            }
        }
    }
</script>