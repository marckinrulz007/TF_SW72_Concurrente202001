<template>
  <v-app>
    <v-app-bar app>
      <v-container class="py-0 fill-height">
        <h1>TRABAJO FINAL</h1>
        <v-spacer></v-spacer>
      </v-container>
    </v-app-bar>
    <v-main>
      <v-container fluid>
        <router-view></router-view>
      </v-container>
    </v-main>
  </v-app>
</template>

<style scoped>
.btn {
  margin-left: 3em;
}
@media screen and (max-width: 480px) {
  .btn{
    margin-left: 1em;
  }
}
</style>
