(function(){
  const API_BASE_DEFAULT = '/';

  function apiBase(){
    const v = document.getElementById('apiBase').value.trim();
    if(!v) return API_BASE_DEFAULT;
    return v.endsWith('/') ? v : v + '/';
  }

  function $(id){return document.getElementById(id)}
  const msg = $('message');
  const attrSection = $('attributesSection');
  const list = $('attributesList');

  function showMessage(t){msg.textContent = t}
  function clearMessage(){msg.textContent = ''}

  async function loadAttributes(){
    clearMessage();
    list.innerHTML = '';
    const personId = $('personId').value.trim();
    if(!personId){showMessage('Enter a person ID');return}
    try{
      const res = await fetch(apiBase() + 'persons/' + encodeURIComponent(personId) + '/attributes');
      if(!res.ok){ showMessage('Failed to load attributes: ' + res.status); attrSection.classList.remove('hidden'); return }
      const data = await res.json();
      renderAttributes(data);
      attrSection.classList.remove('hidden');
      showMessage('Loaded ' + data.length + ' attribute(s)');
    }catch(e){showMessage('Error: ' + e.message); attrSection.classList.remove('hidden')}
  }

  function renderAttributes(attrs){
    list.innerHTML = '';
    if(!Array.isArray(attrs)) return;
    attrs.forEach(a => {
      const li = document.createElement('li'); li.className = 'attribute';
      const keyIn = document.createElement('input'); keyIn.type = 'text'; keyIn.value = a.key || '';
      const valIn = document.createElement('input'); valIn.type = 'text'; valIn.value = a.value || '';
      const info = document.createElement('div'); info.className = 'meta'; info.textContent = 'id: ' + (a.id||'') + (a.updatedAt? ' • updated: '+a.updatedAt : '');
      const saveBtn = document.createElement('button'); saveBtn.textContent = 'Save';
      const delBtn = document.createElement('button'); delBtn.textContent = 'Delete';

      saveBtn.addEventListener('click', async ()=>{
        await updateAttribute(a.id, keyIn.value, valIn.value);
      });
      delBtn.addEventListener('click', async ()=>{
        if(!confirm('Delete attribute '+a.key+'?')) return;
        await deleteAttribute(a.id);
      });

      li.appendChild(keyIn); li.appendChild(valIn); li.appendChild(saveBtn); li.appendChild(delBtn); li.appendChild(info);
      list.appendChild(li);
    })
  }

  function metaForRequest(){
    return { caller: 'webui', reason: 'user-action', traceId: String(Date.now()) };
  }

  async function createAttribute(key, value){
    clearMessage();
    const personId = $('personId').value.trim();
    const body = { key, value, meta: metaForRequest() };
    try{
      const res = await fetch(apiBase() + 'persons/' + encodeURIComponent(personId) + '/attributes', {
        method: 'POST', headers: {'Content-Type':'application/json'}, body: JSON.stringify(body)
      });
      if(res.status === 201){ showMessage('Attribute created'); await loadAttributes(); return }
      const txt = await res.text(); showMessage('Create failed: ' + res.status + ' ' + txt);
    }catch(e){ showMessage('Create error: ' + e.message) }
  }

  async function updateAttribute(attributeId, key, value){
    clearMessage();
    const personId = $('personId').value.trim();
    const body = { key, value, meta: metaForRequest() };
    try{
      const res = await fetch(apiBase() + 'persons/' + encodeURIComponent(personId) + '/attributes/' + attributeId, {
        method: 'PUT', headers: {'Content-Type':'application/json'}, body: JSON.stringify(body)
      });
      if(res.ok){ showMessage('Attribute updated'); await loadAttributes(); return }
      const txt = await res.text(); showMessage('Update failed: ' + res.status + ' ' + txt);
    }catch(e){ showMessage('Update error: ' + e.message) }
  }

  async function deleteAttribute(attributeId){
    clearMessage();
    const personId = $('personId').value.trim();
    try{
      const res = await fetch(apiBase() + 'persons/' + encodeURIComponent(personId) + '/attributes/' + attributeId, { method: 'DELETE' });
      if(res.ok){ showMessage('Attribute deleted'); await loadAttributes(); return }
      const txt = await res.text(); showMessage('Delete failed: ' + res.status + ' ' + txt);
    }catch(e){ showMessage('Delete error: ' + e.message) }
  }

  document.getElementById('loadBtn').addEventListener('click', loadAttributes);
  document.getElementById('createForm').addEventListener('submit', function(e){
    e.preventDefault();
    const k = $('newKey').value.trim(); const v = $('newValue').value;
    if(!k){ showMessage('Key required'); return }
    createAttribute(k,v);
    this.reset();
  });

  // Exported for convenience when opening devtools
  window.personAttributesUI = { loadAttributes };
})();
