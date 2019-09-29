<template>
    <v-container>
        <v-row fill-height style="text-align:center">
            <v-col
            cols="12"
            >
                <h1>BATTLECARDS</h1>
                <h2> Collect your own star wars trading cards </h2>
            </v-col>
        </v-row>
        <v-row :justify="alignment">
            <v-col cols="6">
                <v-text-field
                filled
                v-model="character"
                background-color="white"
                label="Search for a Character"
                append-icon="fas fa-search"
                class = "px-2"
                @click:append="search"
                ></v-text-field>
            </v-col>
        </v-row>
        <v-row :justify="alignment">
            <v-col cols="12" class="text-center">
                <h1> YOUR COLLECTION </h1>
                <h3> Last Search: {{lastDateSearch}} ({{secondsAgo}} seconds ago)</h3>
            </v-col>

        </v-row>
        <v-row :justify="alignment">
            <v-col
            v-for ="character in characters"
            :key = "character.id"
            cols=3>
                <v-card 
                class="text-center"
                :hover="true"
                :elevation="4"
                :elevated="true"
                @click="displayCharacter(character)"
                >
                    <p class = "characterName mb-0"> {{character.Name}}</p>
                    <p class ="characterHome"> {{character.HomeWorld.Name}} </p>
                </v-card>
            </v-col>
        </v-row>
        <v-dialog v-model="error" max-width="400">
            <error :message="errorMessage" v-on:closeerror="closeError"></error>
        </v-dialog>
        <v-dialog v-model="characterDialog" persistent max-width="600">
            <character :info="selectedCharacter" v-on:closecharacter="closeCharacter"></character>
        </v-dialog>
    </v-container> 
</template>
<script>
    const axios = require('axios')
    import character from "../components/character"
    import error from "../components/error"
    export default {
        name: "Home",
        components: {
            character,
            error
        },
        data() {
            return {
                character: null,
                characters: [],
                characterNames: [],
                characterDialog: false,
                selectedCharacter: null,
                lastDateSearch: null,
                secondsAgo: null,
                alignment: "center",
                created: true,
                errorMessage: null,
                error: null,
            }
        },
        async created() {
            //immediately get the characters returned from the getCache function
           this.getCache();
        },
        async beforeUpdate(){
            
        },
        methods: {
            async search(){
                axios.post("http://localhost:5555/api/character", {
                    name:this.character
                }).then((response) => {
                    //response.data is the response from our golang server
                    //Check for the status
                    var resp = response.data
                    if(resp["status"] == false) {
                        this.errorMessage = resp["message"]
                        this.error = true
                    }else{
                        //first check if the character array already contains the characters to be searched
                        if(this.characterNames.includes(resp["message"].Name)) {
                            this.errorMessage = "You've already added this character to your collection!"
                            this.error = true
                        }else{
                            this.characterNames.push(resp["message"].Name)
                            this.characters.push(resp["message"])
                            this.updateSearchDate()
                        }
                    }
                })
            },
            async displayCharacter(character) {
                this.selectedCharacter = character
                this.characterDialog = true
            },            
            async updateSearchDate(){
                this.lastDateSearch = Math.round((new Date()).getTime() / 1000);
                this.updateSearchSecond()
            },
            async updateSearchSecond(){
                var currentTime = Math.round((new Date()).getTime() / 1000)
                this.secondsAgo = currentTime - this.lastDateSearch
            },
            //getCache will get the cache from the backend, and serve the data located in the cache
            async getCache() {
                axios.get("http://localhost:5555/api/getCache").then((response) => {
                    const resp = response.data
                    this.characters = resp["message"].Characters
                    this.lastDateSearch = resp["message"].Time
                    for(var i =0;i<this.characters.length;i++){
                        this.characterNames.push(this.characters[i].Name)
                    }
                    this.updateSearchSecond()
                })
            },
            async closeError() {
                this.error = false
            },
            
            async closeCharacter(){
                this.characterDialog = false
            }
        }
    }

</script>
<style>
   .characterName {
       font-size: 20px
   }
   .characterHome {
       font-size: 15px
   }
</style>