<template>
  <div id="registrationWrapper">
    <div id="mask"></div>
    <iframe name="iframe1" style="position: absolute; left: -9999px;"></iframe>
    <div id="registrationBlock" ref="registrationBlock">

      <app-control v-bind:head="'Регистрация'" v-bind:move="false" v-bind:close="false"
                   v-bind:refWindow="'registrationBlock'"/>

      <form method="POST" id="formNewUser" autocomplete="off" onsubmit="return false;">
        <table>
          <tr>
            <td> Логин:</td>
            <td><input type="text" autocomplete="new-password" v-model="username" placeholder="Login"></td>
          </tr>
<!--          <tr>-->
<!--            <td> Почта:</td>-->
<!--            <td><input type="text" autocomplete="new-password" v-model="email" placeholder="E-mail"></td>-->
<!--          </tr>-->
          <tr>
            <td> Пароль:</td>
            <td><input type="password" autocomplete="new-password" v-model="password" placeholder="password"></td>
          </tr>
          <tr>
            <td> Еще раз:</td>
            <td><input type="password" autocomplete="new-password" v-model="confirm" placeholder="password"></td>
          </tr>
          <tr id="buttonCell">
            <td colspan="2" style="text-align: center; padding-top: 6px;">
              <input v-if="!success" class="button" id="regButton" type="submit" value="Регистрация"
                     @click="registration">
              <input v-if="success" class="button" id="regButton" type="submit" value="Войти"
                     style="box-shadow: 0 0 10px 2px #b4ffb4, inset 0 0 10px 2px #b4ffb4"
                     @click="login">
            </td>
          </tr>
        </table>
        <div class="Failed" id="error">{{err}}</div>
      </form>
    </div>
  </div>
</template>

<script>
  import Control from '../Window/Control';
  import {urls} from '../../const'

  export default {
    name: "Registration",
    data() {
      return {
        username: null,
        password: null,
        confirm: null,
        email: null,
        err: '',
        success: false,
      }
    },
    created() {
      try {
        this.$store.getters.getWSConnectState.socket.close()
      } catch (e) {
      }
    },
    methods: {
      registration() {
        let app = this;
        app.err = '';

        this.$api.post(urls.regURL, {
          username: app.username,
          email: '123',
          password: app.password,
          confirm: app.confirm,
        }, {
          withCredentials: true,
        }).then(function (response) {
          if (response.data.success) {
            app.success = true;
          } else {
            if (response.data.error === "form is empty") {
              app.err = "Не все формы заполнены"
            }
            if (response.data.error === "login busy") {
              app.err = "Этот логин уже занят"
            }
            if (response.data.error === "email busy") {
              app.err = "Эта почта уже используется"
            }
            if (response.data.error === "password error") {
              app.err = "Пароли не совпадают"
            }
            //console.log(response)
          }
        }).catch(function (error) {
          //console.log(error)
        });
      },
      login: function () {

        let app = this;
        app.err = '';

        this.$api.post(urls.loginURL, {
          username: this.username,
          password: this.password,
        }, {
          withCredentials: true,
        }).then(function (response) {
          if (response.data.success) {
            app.load();
          } else {
            if (response.data.error === "not allow") {
              app.err = "Не верный логин либо пароль"
            }
            //console.log(response)
          }
        }).catch(function (error) {
          //console.log(error);
        });
      },
      load() {
        this.$router.push({path: '/gate'});
      }
    },
    components: {
      AppControl: Control,
    }
  }
</script>

<style scoped>
  #registrationWrapper {
    height: 100vh;
    width: 100%;
    text-align: center;
    background-color: #7f7f7f;
    background-image: url('../../assets/bases/base.jpg');
    background-size: cover;
    background-attachment: fixed;
    background-position: center;
  }

  #registrationBlock {
    position: absolute;
    left: calc(50% - 130px);
    top: 20%;
    display: block;
    border-radius: 5px;
    width: 260px;
    height: 160px;
    border: 1px solid #25a0e1;
    background: rgb(8, 138, 210);
    z-index: 11;
    padding: 20px 0 0 0;
    box-shadow: 0 0 2px black;
  }

  #formNewUser input[type="text"], #formNewUser input[type="password"] {
    text-shadow: none;
    color: black;
  }

  #formNewUser input[type="submit"] {
    width: 74px;
  }

  #buttonCell {
    text-align: center
  }

  td {
    width: 150px;
    height: 20px;
    padding: 10px;
    text-shadow: 0 -1px 1px #000000, 0 -1px 1px #000000, 0 1px 1px #000000, 0 1px 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000;
    text-align: right;
    font-size: 12px;
    color: white;
  }

  th {
    text-shadow: 0 1px 1px rgba(0, 0, 0, .3);
  }

  #mask {
    left: 0;
    top: 0;
    background: rgba(0, 0, 0, 0.7);
    position: absolute;
    z-index: 1;
    height: 100%;
    width: 100%;
  }
</style>
