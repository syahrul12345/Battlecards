<template>
    <v-card>
        <v-card-title> {{info.Name}} </v-card-title>
        <v-card-text class="pb-0"> Gender: {{info.Gender}}</v-card-text>
        <v-card-text class ="pt-1 pb-0"> Homeworld: {{info.HomeWorld.Name}} </v-card-text>
        <v-card-text class ="pt-1 pb-0"> Climate: {{info.HomeWorld.Climate}} </v-card-text>
        <v-card-text class ="pt-1 pb-0"> Population: {{info.HomeWorld.Population}} </v-card-text>
        <v-container>
            <v-row>
                <v-col v-if="info.Vehicles != null" cols ="12">
                    <h1 class="text-center"> Vehicles </h1>
                    <vehicle 
                    v-for="vehicle in info.Vehicles"
                    v-bind:key="vehicle.Name"
                    :vehicle ="vehicle"
                    ></vehicle>
                </v-col>
                <v-col v-else cols="12">
                    <h1 class="text-center"> This character has no vehicles </h1>
                </v-col>
            </v-row>
            <v-row>
                <v-col v-if="info.StarShips != null" cols="12">
                    <h1 class="text-center"> StarShips </h1>
                    <starship
                    v-for="starship in info.StarShips"
                    v-bind:key="starship.Model"
                    :starship="starship"></starship>
                </v-col>
                <v-col v-else cols="12">
                    <h1 class="text-center"> This character has no starships </h1>
                </v-col>
            </v-row>
        </v-container>
        <v-card-actions>
            <v-btn @click="closeCharacterDialog">Close</v-btn>
        </v-card-actions>
    </v-card>
</template>
<script>
import vehicle from "./vehicle"
import starship from "./starship"

export default {
    props:['info'],
    components: {
        vehicle,
        starship
    },
    data() {
        return {
            created:true
        }
    },
    methods: {
        async closeCharacterDialog() {
            this.$emit('closecharacter')
        }
    }
}
</script>
<style></style>       