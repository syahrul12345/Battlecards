<template>
    <v-card class="my-1 mx-2">
        <v-card-title> {{vehicle.Name}} </v-card-title>
        <v-card-text class="pb-0"> Model : {{vehicle.Model}}</v-card-text>
        <v-card-text class="pt-1"> Cost : {{vehicle.CostInCredits}} </v-card-text>
    </v-card>
</template>
<script>

export default {
    props:['vehicle']
}

</script>
<style></style>    