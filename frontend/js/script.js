document.addEventListener('DOMContentLoaded', () => {
    const baseURL = 'http://localhost:3030';

    const addForm = document.getElementById('addForm');
    const itemInput = document.getElementById('itemInput');
    const addItemBtn = document.getElementById('addItemBtn');
    const itemList = document.getElementById('itemList');
    const itemCount = document.getElementById('itemCount');

    const parseResponse = async (res) => {
        try { return await res.json(); } catch { return await res.text(); }
    };

    const handleFetch = async (url, opts = {}) => {
        const res = await fetch(url, opts);
        const body = await parseResponse(res);
        if (!res.ok) throw { status: res.status, body };
        return body;
    };

    const renderList = (items) => {
        if (!itemList) return;
        itemList.innerHTML = '';
        items.forEach(it => {
            const li = document.createElement('li');
            li.className = 'flex items-center justify-between p-4 bg-white border border-slate-100 rounded-xl hover:shadow-md transition-shadow';

            const txt = document.createElement('span');
            txt.className = 'text-slate-700 font-medium';
            txt.textContent = it.content;

            const controls = document.createElement('div');
            controls.className = 'flex gap-2';

            const editBtn = document.createElement('button');
            editBtn.className = 'text-blue-500 hover:bg-blue-50 p-2 rounded-lg transition-colors';
            editBtn.textContent = 'Edit';
            editBtn.addEventListener('click', async () => {
                const newVal = prompt('Edit item', it.content);
                if (newVal == null) return;
                const trimmed = newVal.trim();
                if (!trimmed || trimmed === it.content) return;
                try {
                    await updateContent(it.content, trimmed);
                } catch (e) {
                    console.error('update error', e);
                }
            });

            const delBtn = document.createElement('button');
            delBtn.className = 'text-red-500 hover:bg-red-50 p-2 rounded-lg transition-colors';
            delBtn.textContent = 'Delete';
            delBtn.addEventListener('click', async () => {
                if (!confirm('Delete this item?')) return;
                try {
                    await deleteContent(it.content);
                } catch (e) {
                    console.error('delete error', e);
                }
            });

            controls.appendChild(editBtn);
            controls.appendChild(delBtn);

            li.appendChild(txt);
            li.appendChild(controls);
            itemList.appendChild(li);
        });
        if (itemCount) itemCount.textContent = String(items.length);
    };

    const getAll = async () => {
        try {
            const data = await handleFetch(`${baseURL}/get`);
            renderList(Array.isArray(data) ? data : []);
            await getTotal(); // keeps server count in sync if needed
        } catch (e) {
            console.error('getAll error', e);
        }
    };

    const postContent = async (content) => {
        return handleFetch(`${baseURL}/post`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ content })
        });
    };

    const deleteContent = async (content) => {
        await handleFetch(`${baseURL}/delete`, {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ content })
        });
        await getAll();
    };

    const updateContent = async (oldContent, newContent) => {
        await handleFetch(`${baseURL}/update`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ old_content: oldContent, new_content: newContent })
        });
        await getAll();
    };

    const getTotal = async () => {
        try {
            const res = await handleFetch(`${baseURL}/totalnum`);
            if (itemCount && res && typeof res.total === 'number') {
                itemCount.textContent = String(res.total);
            }
        } catch (e) {
            console.error('total error', e);
        }
    };

    if (addForm) {
        addForm.addEventListener('submit', async (ev) => {
            ev.preventDefault();
            const val = itemInput ? itemInput.value.trim() : '';
            if (!val) return;
            try {
                await postContent(val);
                if (itemInput) itemInput.value = '';
                await getAll();
            } catch (e) {
                console.error('post error', e);
            }
        });
    } else if (addItemBtn) {
        addItemBtn.addEventListener('click', (ev) => ev.preventDefault());
    }

    // initial load
    getAll();
});