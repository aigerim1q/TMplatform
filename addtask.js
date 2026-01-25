const modal = document.getElementById('modal');
const overlay = document.getElementById('overlay');
const userList = document.getElementById('userList');
const searchInput = document.getElementById('searchInput');
const checkbox = document.getElementById('executorCheckbox');
const field = document.getElementById('executorField');

const users = [
  { name: 'Arman', role: 'Mid Front End Developer' },
  { name: 'Arailyn', role: 'UX / UI designer' },
  { name: 'Algerim', role: 'Product manager' }
];

let selectedUser = null;

/* OPEN */
function openModal() {
  modal.style.display = 'block';
  overlay.style.display = 'block';
  render(users);
}

/* CLOSE */
function closeModal() {
  modal.style.display = 'none';
  overlay.style.display = 'none';

  if (!selectedUser) {
    field.classList.add('error');
    checkbox.checked = false;
  }
}

/* CONFIRM */
function confirmSelection() {
  if (!selectedUser) {
    field.classList.add('error');
    return;
  }

  checkbox.checked = true;
  field.classList.remove('error');
  closeModal();
}

/* SEARCH */
searchInput.addEventListener('input', e => {
  const value = e.target.value.toLowerCase();
  render(users.filter(u => u.name.toLowerCase().includes(value)));
});

/* RENDER */
function render(list) {
  userList.innerHTML = '';
  list.forEach(user => {
    const div = document.createElement('div');
    div.className = 'user' + (selectedUser === user ? ' active' : '');
    div.innerHTML = `
      <div>
        <strong>${user.name}</strong><br>
        <small>${user.role}</small>
      </div>
      <div class="check"></div>
    `;
    div.onclick = () => {
      selectedUser = user;
      render(list);
    };
    userList.appendChild(div);
  });
}

/* OVERLAY CLICK */
overlay.onclick = closeModal;
