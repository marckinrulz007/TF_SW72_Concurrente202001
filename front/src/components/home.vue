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
        { text: "Puerto", value: "port" },
        { text: "Hash", value: "hash" },
        { text: "Hash previo", value: "previous_hash" },
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
              port: item[0],
              hash: item[1],
              previous_hash: item[2],
            });
          }
          this.displayData = result;
        })
        .catch((e) => {
          console.log(e);
        });
    },
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
