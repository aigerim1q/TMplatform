const startCallBtn = document.getElementById('startCallBtn');
const closeModalBtn = document.getElementById('closeModalBtn');
const inviteBtn = document.getElementById('inviteBtn');
const modal = document.getElementById('callModal');
const jitsiContainer = document.getElementById('jitsiContainer');

let api = null;
let roomName = null;

function generateRoomName() {
  return 'project-call-' + Math.random().toString(36).substr(2, 9);
}

startCallBtn.addEventListener('click', () => {
  modal.classList.add('open');

  if (!roomName) {
    roomName = generateRoomName();
  }

  api = new JitsiMeetExternalAPI('meet.jit.si', {
    roomName,
    parentNode: jitsiContainer,
    width: '100%',
    height: '100%',
    configOverwrite: {
      startWithAudioMuted: true,
      startWithVideoMuted: false,
    },
  });
});

closeModalBtn.addEventListener('click', () => {
  modal.classList.remove('open');

  if (api) {
    api.dispose();
    api = null;
  }
});

inviteBtn.addEventListener('click', async () => {
  const link = `https://meet.jit.si/${roomName}`;
  await navigator.clipboard.writeText(link);
  alert('Ссылка на звонок скопирована');
});
