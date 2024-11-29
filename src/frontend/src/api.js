const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;


// API WIKIS

export async function getAllWikis() {
  const resp = await fetch(`${API_BASE_URL}/wikis`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  if (!resp.ok) {
    throw new Error("Failed to fetch wikis");
  }
  return resp.json();
}

export async function postWiki(data) {
  const resp = await fetch(`${API_BASE_URL}/wikis`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  if (!resp.ok) {
    throw new Error("Failed to post wiki");
  }
  return resp.json();
}

export async function getWiki(id) {
  const resp = await fetch(`${API_BASE_URL}/wikis/${id}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  if (!resp.ok) {
    throw new Error("Failed to fetch wikis");
  }
  return resp.json();
}

export async function deleteWiki(id) {
  const resp = await fetch(`${API_BASE_URL}/wikis/${id}`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  if (!resp.ok) {
    throw new Error("Failed to delete wiki with id " + id);
  }
  return resp.json();
}

export async function putWiki(id){
  const resp = await fetch(`${API_BASE_URL}/wikis/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  if (!resp.ok) {
    throw new Error("Failed to fetch wiki with id " + id);
  }
  return resp.json();
}

export async function getWikiByExactTitle(title) {
  const resp = await fetch(`${API_BASE_URL}/wikis/exact_title?title=${title}`);
  if (!resp.ok) {
    throw new Error("Failed to fetch wikis");
  }
  return resp.json();
}

export async function getWikiByTitle(title) {
  const resp = await fetch(`${API_BASE_URL}/wikis/title?title=${title}`);
  if (!resp.ok) {
    throw new Error("Failed to fetch wikis");
  }
  return resp.json();
}

export async function getWikiByDescription(description) {
  const resp = await fetch(`${API_BASE_URL}/wikis/description?description=${description}`);
  if (!resp.ok) {
    throw new Error("Failed to fetch wikis");
  }
  return resp.json();
}

export async function getWikiByCategory(category) {
  const resp = await fetch(`${API_BASE_URL}/wikis/category?category=${category}`);
  if (!resp.ok) {
    throw new Error("Failed to fetch wikis");
  }
  return resp.json();
}

// API ENTRADAS

//get all entries
export async function getAllEntries() {
  const resp = await fetch(`${API_BASE_URL}/entries`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  if (!resp.ok) {
    throw new Error("Failed to fetch entries");
  }
  return resp.json();
}

//post an entry

export async function postEntry(data) {
  const resp = await fetch(`${API_BASE_URL}/entries`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  if (!resp.ok) {
    throw new Error("Failed to post entry");
  }
  return resp.json();
}

//get an entry by id

export async function getEntry(id) {
  const resp = await fetch(`${API_BASE_URL}/entries/${id}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  if (!resp.ok) {
    throw new Error("Failed to fetch entry");
  }
  return resp.json();
}

//delete an entry by id

export async function deleteEntry(id) {
  const resp = await fetch(`${API_BASE_URL}/entries/${id}`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  if (!resp.ok) {
    throw new Error("Failed to delete entry with id " + id);
  }
  return resp.json();
}

//put an entry by id

export async function putEntry(id){
  const resp = await fetch(`${API_BASE_URL}/entries/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  if (!resp.ok) {
    throw new Error("Failed to fetch entry with id " + id);
  }
  return resp.json();
}

//get an entry by wiki

export async function getEntryByWikiId(wikiId) {
  const resp = await fetch(`${API_BASE_URL}/entries/wiki_id?=${wikiId}`);
  if (!resp.ok) {
    throw new Error("Failed to fetch entries");
  }
  return resp.json();
}

//get an entry by exact title

export async function getEntryByExactTitle(title) {
  const resp = await fetch(`${API_BASE_URL}/entries/exact_title?title=${title}`);
  if (!resp.ok) {
    throw new Error("Failed to fetch entries");
  }
  return resp.json();
}

//get an entry by title

export async function getEntryByTitle(title) {
  const resp = await fetch(`${API_BASE_URL}/entries/title?title=${title}`);
  if (!resp.ok) {
    throw new Error("Failed to fetch entries");
  }
  return resp.json();
}

//get an entry by author

export async function getEntryByAuthor(author) {
  const resp = await fetch(`${API_BASE_URL}/entries/author?author=${author}`);
  if (!resp.ok) {
    throw new Error("Failed to fetch entries");
  }
  return resp.json();
}

//get an entry by date

export async function getEntryByDate(date) {
  const resp = await fetch(`${API_BASE_URL}/entries/date?date=${date}`);
  if (!resp.ok) {
    throw new Error("Failed to fetch entries");
  }
  return resp.json();
}


// API VERSIONES

//get all versions

export async function getAllVersions() {
  const resp = await fetch(`${API_BASE_URL}/versions`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  if (!resp.ok) {
    throw new Error("Failed to fetch versions");
  }
  return resp.json();
}

//post a version

export async function postVersion(data) {
  const resp = await fetch(`${API_BASE_URL}/versions`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  if (!resp.ok) {
    throw new Error("Failed to post version");
  }
  return resp.json();
}

//get a version by id

export async function getVersion(id) {
  const resp = await fetch(`${API_BASE_URL}/versions/${id}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  if (!resp.ok) {
    throw new Error("Failed to fetch version");
  }
  return resp.json();
}

//delete a version by id

export async function deleteVersion(id) {

  const resp = await fetch(`${API_BASE_URL}/versions/${id}`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!resp.ok) {
    throw new Error("Failed to delete version with id " + id);
  }
  return resp.json();
}

//put a version by id

export async function putVersion(id){
  const resp = await fetch(`${API_BASE_URL}/versions/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  if (!resp.ok) {
    throw new Error("Failed to fetch version with id " + id);
  }
  return resp.json();
}

//get all versions by entry id

export async function getVersionsByEntryId(entryId) {
  const resp = await fetch(`${API_BASE_URL}/versions/entry?entry=${entryId}`);
  if (!resp.ok) {
    throw new Error("Failed to fetch versions");
  }
  return resp.json();
}


// API COMENTARIOS


// API MEDIA