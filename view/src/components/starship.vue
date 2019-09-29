<template>
    <v-card class="my-1 mx-2">
        <v-card-title> {{starship.Model}} </v-card-title>
        <v-card-text class="pb-0"> Class : {{starship.Starship_Class}}</v-card-text>
        <v-card-text class="pt-1 pb-0"> HyperDrive Rating : {{starship.Hyperdrive_Rating}} </v-card-text>
        <v-card-text class="pt-1 pb-0" > Cost: {{starship.Cost_In_Credits}} </v-card-text>
        <v-card-text class ="pt-1"> Manufacturer: {{starship.Manufacturer}} </v-card-text>
    </v-card>
</template>
<script>

export default {
    props:['starship']
}

</script>
<style></style>    