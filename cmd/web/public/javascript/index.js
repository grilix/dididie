import { createClient } from './client.js';

const didDieToday = (profile) => {
  if (!profile.died_at) {
    return false;
  }

  const now = new Date();
  const diedAt = profile.died_at;

  return now.toLocaleDateString() === diedAt.toLocaleDateString();
};

const q = (selector) => {
  return document.querySelectorAll(selector);
};

const showElements = (elements) => {
  elements.forEach(e => e.classList.remove('hidden'));
};

const showStatus = (profile) => {
  queryForm.elements.id.value = profile.id;

  q('.status').forEach(e => e.classList.add('hidden'));

  if ('id' in profile) {
    if (didDieToday(profile)) {
      showElements(q('.status-died'));
    } else {
      showElements(q('.status-notDied'));
    }
  } else {
    showElements(q('.status-neverDied'));
  }
};

const onSubmitForm = (form, fn) => {
  form.addEventListener('submit', (e) => {
    e.stopPropagation();
    e.preventDefault();

    fn();

    return false;
  });
};

const client = createClient('http://localhost:8080');

const queryForm = document.getElementById('query-form');

onSubmitForm(queryForm, () => {
  client.query(queryForm.elements.name.value).then(showStatus);
});

const dieForm = document.getElementById('die-form');

onSubmitForm(dieForm, () => {
  client.die(queryForm.elements.id.value).then(showStatus);
});

const startForm = document.getElementById('start-form');

onSubmitForm(startForm, () => {
  client.start(queryForm.elements.name.value).then(showStatus);
});
