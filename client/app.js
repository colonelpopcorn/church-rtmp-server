const BASE_URL = `${window.location.origin}/api`;

const fetchWrapper = {
  get,
  post,
  put,
  delete: _delete,
};

const SESSION_KEY = "sessionToken";
const DEFAULT_HEADERS = () => {
  const sessionToken = localStorage.getItem(SESSION_KEY);
  return { Authorization: sessionToken ? `Bearer ${sessionToken}` : "" };
};

function get(url) {
  const requestOptions = {
    method: "GET",
    headers: { ...DEFAULT_HEADERS() },
  };
  return fetch(url, requestOptions).then(handleResponse);
}

function post(url, body) {
  const requestOptions = {
    method: "POST",
    headers: { ...DEFAULT_HEADERS(), "Content-Type": "application/json" },
    body: JSON.stringify(body),
  };
  return fetch(url, requestOptions).then(handleResponse);
}

function put(url, body) {
  const requestOptions = {
    method: "PUT",
    headers: { ...DEFAULT_HEADERS(), "Content-Type": "application/json" },
    body: JSON.stringify(body),
  };
  return fetch(url, requestOptions).then(handleResponse);
}

// prefixed with underscored because delete is a reserved word in javascript
function _delete(url) {
  const requestOptions = {
    method: "DELETE",
    headers: { ...DEFAULT_HEADERS() },
  };
  return fetch(url, requestOptions).then(handleResponse);
}

// helper functions

function handleResponse(response) {
  return response.text().then((text) => {
    const data = text && JSON.parse(text);

    if (!response.ok) {
      const error = (data && data.message) || response.statusText;
      return Promise.reject(error);
    }

    return data;
  });
}

Vue.use(window.VueCodemirror);

const StreamManager = {
  created: function () {
    this.getStreams();
  },
  data: function () {
    return {
      streamsLoaded: false,
      streams: [],
      video: null,
    };
  },
  methods: {
    getStreams() {
      fetchWrapper
        .get(`${BASE_URL}/streams/`)
        .then((data) => {
          this.streams = data.streams;
          this.streamsLoaded = true;
        })
        .catch(console.error);
    },
    createStream() {
      fetchWrapper
        .post(`${BASE_URL}/streams/create-key`)
        .then((_) => this.getStreams());
    },
    deleteStream(id) {
      fetchWrapper
        .delete(`${BASE_URL}/streams/${id}`)
        .then((_) => this.getStreams());
    },
    setStreamValue(streamKey) {
      this.video = document.getElementById("video");
      if (Hls.isSupported()) {
        var hls = new Hls();
        hls.loadSource(`${window.location}live/${streamKey}/index.m3u8`);
        hls.attachMedia(video);
      }
    },
    streamModalClosed() {
      this.video = null;
    },
  },
  template: `
  <div v-else>
    <div v-if="streamsLoaded" id="results-table">
        <table class="table">
            <thead>
                <tr>
                    <th scope="col">Stream Id</th>
                    <th scope="col">Stream Key</th>
                    <th scope="col">Live?</th>
                </tr>
            </thead>
            <tbody v-if="streams.length > 0">
                <tr v-for="stream in streams">
                    <th scope="row">{{ stream.streamId }}</th>
                    <td>{{ stream.streamKey }}</td>
                    <td>
                        <div v-if="stream.isValidStream">
                            <button type="button" class="btn btn-danger" data-toggle="modal" data-target="#videoModal"
                                @click="setStreamValue(stream.streamKey)">
                                LIVE
                            </button>
                        </div>
                    </td>
                    <td>
                        <button type="button" class="btn btn-danger" @click="deleteStream(stream.streamId)">
                            Delete Stream
                        </button>
                    </td>
                </tr>
            </tbody>
            <tbody v-else>
                <tr>
                    <td>No streams yet! Create one and get to work!!</td>
                </tr>
            </tbody>
        </table>
        <!-- Modal -->
        <div class="modal fade" id="videoModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
            aria-hidden="true">
            <div class="modal-dialog modal-dialog-centered modal-lg" role="document">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title" id="exampleModalLongTitle">
                            Stream viewer
                        </h5>
                        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                            <span aria-hidden="true">&times;</span>
                        </button>
                    </div>
                    <div class="modal-body">
                        <video id="video" controls preload="auto" width="768" height="432"></video>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" data-dismiss="modal"
                            @click="streamModalClosed()">
                            Close
                        </button>
                    </div>
                </div>
            </div>
        </div>
        <!-- End Modal -->
        <div>
            <button type="button" class="btn btn-primary" @click="createStream()">
                Create Stream
            </button>
        </div>
    </div>
    <div v-else class="custom-spinner spinner-border" role="status">
        <span class="sr-only">Loading...</span>
    </div>
  </div>
  `,
};

const ConfigEditor = {
  computed: {
    nginxConfUnchanged: (vm) => vm.initialNginxContent === vm.nginxContent,
  },
  created: function () {
    this.openConfEditor();
  },
  data: function () {
    return {
      nginxContent: "",
      initialNginxContent: "",
      confEditorOpen: false,
      codeMirrorOpts: {
        theme: "eclipse",
        lineNumbers: true,
        viewportMargin: Infinity,
        mode: "text/nginx",
      },
    };
  },
  methods: {
    async openConfEditor() {
      this.confEditorOpen = true;
      if (this.nginxContent === "") {
        fetchWrapper.get(`${BASE_URL}/nginx/config`).then((res) => {
          this.initialNginxContent = res.content;
          this.nginxContent = this.initialNginxContent;
        });
      }
    },
    closeConfEditor() {
      this.confEditorOpen = false;
    },
    submitNewContent() {
      fetchWrapper
        .post(`${BASE_URL}/nginx/config`, { content: this.nginxContent })
        .then(() => (this.initialNginxContent = this.nginxContent))
        .catch(console.error);
    },
  },
  template: `
  <div>
    <!-- Conf editor -->
    <div v-if="nginxContent !== ''">
      <codemirror ref="cmEditor" v-model="nginxContent" :options="codeMirrorOpts" />
    </div>
    <div v-else class="custom-spinner spinner-border" role="status">
      <span class="sr-only">Loading...</span>
    </div>
    <button :disabled="nginxConfUnchanged" type="button" class="btn btn-primary float-right mt-3" @click="submitNewContent()">
      Submit
    </button>
  </div>
  `,
};

const LoginPage = {
  template: `
  <div id="login-row" class="row justify-content-center align-items-center">
    <div id="login-column" class="col-md-6">
      <div id="login-box" class="col-md-12">
        <form id="login-form" class="form-signin" @submit.prevent="login()">
          <h3 class="text-center">Login</h3>
          <div class="form-group">
            <label for="username" class="sr-only">Email address</label>
            <input
              v-model="username"
              type="text"
              id="username"
              class="form-control"
              placeholder="Username"
              required=""
              autofocus=""
            />
          </div>
          <div class="form-group">
            <label for="password" class="sr-only">Password</label>
            <input
              v-model="password"
              type="password"
              id="password"
              class="form-control"
              placeholder="Password"
              required=""
            />
          </div>
          <div class="form-group">
            <input type="submit" name="submit" class="btn btn-primary btn-block" value="Submit">
          </div>
          <div class="form-group float-right">
            <label for="remember-me"><span>Remember me</span> <span><input id="remember-me" name="remember-me" type="checkbox"></span></label>
          </div>
        </form>
      </div>
    </div>
  </div>`,
  data() {
    return {
      username: "",
      password: "",
    };
  },
  methods: {
    async login() {
      const response = await fetchWrapper.post(`${BASE_URL}/login`, {
        username: this.username,
        password: this.password,
      });
      if (response) {
        if (response.token) {
          localStorage.setItem(SESSION_KEY, response.token);
        }
        if (response.isAdmin && response.routes) {
          STREAMING_APP.routes = response.routes;
          STREAMING_APP.isAdmin = response.isAdmin;
        }
        this.$router.push("/");
      }
    },
  },
};

const AuthService = {
  async isAuthenticated() {
    const sessionToken = localStorage.getItem(SESSION_KEY);
    return !!sessionToken;
  },
};

const AuthGuard = async (to, _, next) => {
  if (to.meta.requiresAuth) {
    const isAuthenticated = await AuthService.isAuthenticated();
    if (isAuthenticated) {
      next();
    } else {
      next("/login");
    }
  } else {
    next();
  }
};

const routes = [
  { path: "/", redirect: "/stream-manager" },
  { path: "/login", component: LoginPage },
  {
    path: "/stream-manager",
    component: StreamManager,
    meta: { requiresAuth: true },
  },
  {
    path: "/config-editor",
    component: ConfigEditor,
    meta: { requiresAuth: true },
  },
];

const router = new VueRouter({
  routes,
});

router.beforeEach(AuthGuard);

const STREAMING_APP = new Vue({
  router,
  data() {
    return {
      isAdmin: false,
      routes: [],
    };
  },
  methods: {
    logout() {
      fetchWrapper.post(`${BASE_URL}/logout`, {}).then(() => {
        localStorage.setItem(SESSION_KEY, "");
        this.isAdmin = false;
        this.routes = [];
      });
    },
  },
}).$mount("#app");
