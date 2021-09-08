<template>
  <div class="clip-input">
    <br />
    <table class="center">
      <tr>
        <th>url</th>
        <th>
          <input
            type="text"
            id="url"
            placeholder="https://www.youtube.com/watch?v=HJW-2lSM_Sk"
            class="spaced"
            v-model="url"
          />
        </th>
      </tr>
    </table>
    <br />
    <table class="center">
      <tr>
        <th>Klip eleje</th>
        <th>Klip vége</th>
        <th>Játékos</th>
        <th>Szituáció</th>
      </tr>
      <tr v-for="(c, i) in clips" :key="i">
        <td>
          <input
            type="text"
            v-model="c.start"
            name="start"
            placeholder="00:00"
            class="input-text"
          />
        </td>
        <td>
          <input
            type="text"
            v-model="c.end"
            name="start"
            placeholder="02:12"
            class="input-text"
          />
        </td>
        <td>
          <select v-model="c.player" class="select">
            <option v-for="name in names" :key="name" :value="name">
              {{ name }}
            </option>
          </select>
        </td>
        <td>
          <select v-model="c.event" class="select">
            <option v-for="event in events" :key="event" :value="event">
              {{ event }}
            </option>
          </select>
        </td>
      </tr>
      <tr></tr>
      <tr>
        <td colspan="1">
          <input
            type="button"
            value="Új sor"
            class="wide-button"
            @click="addNewRow()"
          />
        </td>
        <td colspan="1">
          <input
            type="button"
            value="Törlés"
            class="wide-button"
            @click="removeRow()"
          />
        </td>
        <td></td>
        <td colspan="1">
          <input
            type="button"
            value="Küld"
            class="wide-button"
            @click="sendClipsToBackend"
          />
        </td>
      </tr>
    </table>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

interface Player {
  fname: string;
  lname: string;
}

// interface Clip {
//   start: string;
//   end: string;
//   player: string;
//   event: string;
// }

export default defineComponent({
  name: "ClipInput",
  created() {
    this.fetchNames();
  },
  data() {
    return {
      clips: [{ start: "", end: "", player: "", event: "" }],
      url: "",
      names: ["Csapat"],
      selectedName: "Csapat",
      events: ["Csel", "Megúszás", "Szerelés", "Szép akció", "Egyéb"],
      selectedEvent: "Csel"
    };
  },
  methods: {
    addNewRow() {
      this.clips.push({
        start: "",
        end: "",
        player: "",
        event: ""
      });
    },
    removeRow() {
      this.clips.pop();
    },
    sendClipsToBackend() {
      const body = {
        url: this.url,
        clips: this.clips
      };
      const options = {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body)
      };

      console.log(options.body);

      // jesus... FIXME: nginx maybe? or :8000/clip works?
      fetch("http://pumate.hu:8000/clip/", options)
        .then(resp => console.log(resp.json()))
        .then(() => (this.clips = []))
        .catch(err => console.log(err));
    },
    fetchNames() {
      fetch("https://www.vizalattihoki.hu/api/v1/players").then(res => {
        res.json().then(data => {
          data.forEach((player: Player) => {
            this.names.push(`${player.lname} ${player.fname}`);
          });
        });
      });
    }
  }
});
</script>

<style scoped>
.center {
  margin-left: auto;
  margin-right: auto;
}

.spaced {
  margin-left: 20px;
  margin-right: 20px;
  width: 400px;
}

.wide-button {
  width: 185px;
}

.input-text {
  text-align: right;
  width: 185px;
}

.select {
  text-align: center;
  width: 185px;
}
</style>
