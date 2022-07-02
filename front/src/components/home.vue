<template>
  <v-card class="card" style="border-radius: 1em">
    <v-card-title class="justify-center">BLOCKCHAIN CONCURRENTE</v-card-title>
    <v-card-subtitle class="text-center pb-8"
      >Información obtenida por la comunicación entre nodos</v-card-subtitle
    >
    <v-data-table :headers="headers" :items="displayData" :items-per-page="15">
    </v-data-table>
  </v-card>
</template>

<script>
import DatasetService from "@/services/dataset-service";

export default {
  name: "home",
  data() {
    return {
      displayData: [],
      headers: [
        { text: "Id", value: "id" },
        { text: "Hash", value: "hash" },
        { text: "Hash previo", value: "previous_hash" },
        { text: "Tiempo", value: "time" },
      ],
    };
  },
  methods: {
    getAllData() {
      DatasetService.getAllData()
        .then((response) => {
          const result = [];
          const data = response.data;
          for (const item of data) {
            result.push({
              id: item[0],
              hash: item[1],
              previous_hash: item[2],
              time: this.fixDates(item[3]),
            });
          }
          this.displayData = result;
        })
        .catch((e) => {
          console.log(e);
        });
    },
    fixDates(string){
      const parts = string.split(" ");
      const date = parts[0].split("-");
      const time = parts[1];
      return `${date[2]}/${date[1]}/${date[0]} ${time.slice(0,5)}`;
    }
  },
  mounted() {
    this.getAllData();
  },
};
</script>

<style scoped>
.card {
  margin-top: 2em;
}
@media screen and (max-width: 700px) {
  .card {
    margin-top: 0;
  }
}
</style>
