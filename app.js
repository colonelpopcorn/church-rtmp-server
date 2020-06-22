const BASE_URL = `${window.location}api`;

const fetchWrapper = {
  get,
  post,
  put,
  delete: _delete,
};

function get(url) {
  const requestOptions = {
    method: 'GET',
  };
  return fetch(url, requestOptions).then(handleResponse);
}

function post(url, body) {
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  };
  return fetch(url, requestOptions).then(handleResponse);
}

function put(url, body) {
  const requestOptions = {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  };
  return fetch(url, requestOptions).then(handleResponse);
}

// prefixed with underscored because delete is a reserved word in javascript
function _delete(url) {
  const requestOptions = {
    method: 'DELETE',
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

var app = new Vue({
  el: '#app',
  data: {
    message: 'Hello Vue!',
    streamsLoaded: false,
    streams: [],
  },
  created: function () {
    this.getStreams();
  },
  methods: {
    getStreams() {
      fetchWrapper
        .get(`${BASE_URL}/streams`)
        .then((data) => {
          this.streams = data.streams;
          this.streamsLoaded = true;
        })
        .catch(console.error);
    },
    createStream() {
      fetchWrapper
        .post(`${BASE_URL}/create-key`)
        .then((_) => this.getStreams());
    },
    deleteStream(id) {
      fetchWrapper
        .delete(`${BASE_URL}/streams/${id}`)
        .then((_) => this.getStreams());
    },
  },
});
