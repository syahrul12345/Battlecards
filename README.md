### Battlecards
Easily search for your a Starwars character, and get their collectible card. Collect as much Starwars character cards as you want!
You can search by name or planets. It also has an analytics tool to compare the stats of all characters in your collection!

## Development
Instructions on how to build and play with your own local copy of Battlecards.

# Prerequisites
You'll need these dependencies before getting started:
```
javascript
nodejs
npm
golang
```

# Installing

Follow these instructions to get Battlecards running after you've installed the dependencies

First, clone the repository and install all required node modules.
```
git clone https://github.com/syahrul12345/Battlecards.git
cd Battlecards
npm install
```

Next, let's compile the frontend.
```
cd view
npm run build
```
The frontend will live in the root folder in a newly generated `/dist/` directory. Make sure `/dist/` is in the root folder. Our backend server will target the `index.html` file and serve it on `http://localhost:5555`.

The server is written in golang and lives in the `backend` folder. It acts as middlelayer between the frontend and the SWAPI, allowing us to re-index the api and to apply caching features. Caching is done by saving it to the `cache` folder in the root directory.

```
cd backend
go build
./backend/ 
```
This generates a binary called `backend`. You can execute the binary using the command

```
./backend/
```

That's it! If you've followed this steps in order, the backend will automatically serve the `index.html` file in the `/dist/` folder. You can visit your own Battlecards website 

# Backend architecture & design

The entry point to the backend is called `main.go`. This stores the logic for all API calls.