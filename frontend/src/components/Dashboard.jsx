import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api';

function Dashboard() {
  const [notes, setNotes] = useState([]);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    fetchNotes();
  }, []);

  const fetchNotes = async () => {
    try {
      const res = await api.get('/notes');
      setNotes(res.data);
    } catch (err) {
      if (err.response && err.response.status === 401) {
        localStorage.removeItem('token');
        navigate('/login');
      }
    }
  };

  const handleCreate = async (e) => {
    e.preventDefault();
    try {
      await api.post('/notes', { title, content });
      setTitle('');
      setContent('');
      fetchNotes();
    } catch (err) {
      setError('Failed to create note');
    }
  };

  const handleDelete = async (id) => {
    try {
      await api.delete(`/notes/${id}`);
      fetchNotes();
    } catch (err) {
      setError('Failed to delete note');
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  return (
    <div className="container">
      <div className="header">
        <h2>Your Notes</h2>
        <button className="btn btn-danger btn-small" onClick={handleLogout}>Logout</button>
      </div>
      
      {error && <div className="error-msg">{error}</div>}

      <div className="glass-panel create-note">
        <h3>Create New Note</h3>
        <form onSubmit={handleCreate}>
          <div className="form-group">
            <input 
              type="text" 
              placeholder="Note Title" 
              value={title} 
              onChange={e => setTitle(e.target.value)} 
              required 
            />
          </div>
          <div className="form-group">
            <textarea 
              placeholder="Write your note here..." 
              value={content} 
              onChange={e => setContent(e.target.value)} 
              required 
            />
          </div>
          <button type="submit" className="btn">Add Note</button>
        </form>
      </div>

      <div className="notes-grid">
        {notes && notes.length > 0 ? (
          notes.map(note => (
            <div key={note.id} className="note-card">
              <h3>{note.title}</h3>
              <p>{note.content}</p>
              <div className="note-actions">
                <button 
                  className="btn btn-danger btn-small" 
                  onClick={() => handleDelete(note.id)}
                >
                  Delete
                </button>
              </div>
            </div>
          ))
        ) : (
          <p style={{color: 'var(--text-muted)'}}>No notes found. Create your first one!</p>
        )}
      </div>
    </div>
  );
}

export default Dashboard;
