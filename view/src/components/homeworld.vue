<template>
    <v-card>
        <v-card-title> {{homeworld.Name}} </v-card-title>
        <v-card-text class="py-0"> {{homeworld.Population}}</v-card-text>
        <v-card-text class="pt-0 pb-1"> {{homeworld.Climate}} </v-card-text>
    </v-card>
</template>
<script>
export default {
    props:['homeworld'],
    components:{

    },
    data() {
        return {
            created: true
        }
    }
}
</script>
<style></style>    