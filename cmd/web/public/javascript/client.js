const decodeQueryResponse = (r) => {
  const { profile } = r;

  if (!profile) {
    return {};
  }

  let died_at;

  if (profile.died_at) {
    died_at = new Date(profile.died_at);
  }

  return {
    ...profile,
    died_at,
  };
};

const die = (url, id) => {
  return fetch(`${url}/profiles/${id}/die`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json;charset=utf8',
    },
  })
    .then((r) => {
      return r.json().then(decodeQueryResponse);
    });
};

const query = (url, name) => {
  return fetch(`${url}/profiles/query`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json;charset=utf8',
    },
    body: JSON.stringify({ name }),
  })
    .then((r) => {
      return r.json().then(decodeQueryResponse);
    });
};

const start = (url, name) => {
  return fetch(`${url}/profiles`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json;charset=utf8',
    },
    body: JSON.stringify({ name }),
  })
    .then((r) => {
      return r.json().then(decodeQueryResponse);
    });
};

export const createClient = (url) => {
  return {
    query: name => query(url, name),
    die: id => die(url, id),
    start: name => start(url, name),
  };
};
