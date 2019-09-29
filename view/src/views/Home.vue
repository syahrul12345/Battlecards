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
                <h3> Last Search: {{lastDateSearch}} </h3>
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
                alignment: "center",
                created: true,
                errorMessage: null,
                error: null,
            }
        },
        async created() {
            //immediately get the characters returned from the getCache function
            var characters = this.getCache();
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
                            this.cache(resp["message"].Name)
                        }
                    }
                })
            },
            async displayCharacter(character) {
                this.selectedCharacter = character
                this.characterDialog = true
            },            
            async updateSearchDate(){
                var currentdate = new Date();
                this.lastDateSearch = currentdate.getDate() 
                + "/"+ (currentdate.getMonth()+1)  + "/" 
                + currentdate.getFullYear()+ " @ "  
                + currentdate.getHours() + ":"  
                + currentdate.getMinutes() + ":" 
                + currentdate.getSeconds();
            },
            //getCache will get the cache from the backend, and serve the data located in the cache
            async getCache() {
                console.log('getCached called')
            },
            //Caches the data after the search is completed
            //calls the /api/cache/ with the current character infroamtion  as a payload
            //server will store it in a testfile
            async cache(charName) {
                console.log('Caching for character ' +charName)
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