const { ipcRenderer } = require('electron');
const { marked } = require('marked');
const hljs = require('highlight.js');

// Initialize store for settings and history
let store;
try {
  const Store = require('electron-store');
  store = new Store();
} catch (error) {
  console.warn('electron-store not available, using localStorage fallback');
  // Fallback to localStorage
  store = {
    get: (key, defaultValue) => {
      const value = localStorage.getItem(key);
      return value ? JSON.parse(value) : defaultValue;
    },
    set: (key, value) => {
      localStorage.setItem(key, JSON.stringify(value));
    }
  };
}


// Configure marked with highlight.js
marked.setOptions({
  highlight: function(code, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(code, { language: lang }).value;
      } catch (__) {}
    }
    return hljs.highlightAuto(code).value;
  },
  breaks: true,
  gfm: true
});

// DOM elements
const welcomeScreen = document.getElementById('welcome-screen');
const contentView = document.getElementById('content-view');
const markdownContent = document.getElementById('markdown-content');
const currentFileSpan = document.getElementById('current-file');
const themeToggle = document.getElementById('theme-toggle');
const themeIcon = document.querySelector('.theme-icon');
const openFileBtn = document.getElementById('open-file');
const welcomeOpenBtn = document.getElementById('welcome-open');
const dropOverlay = document.getElementById('drop-overlay');
const recentFilesList = document.getElementById('recent-files-list');
const highlightTheme = document.getElementById('highlight-theme');
const zoomInBtn = document.getElementById('zoom-in');
const zoomOutBtn = document.getElementById('zoom-out');
const zoomResetBtn = document.getElementById('zoom-reset');
const zoomLevelSpan = document.getElementById('zoom-level');

// State
let currentFilePath = null;
let isDarkTheme = store.get('darkTheme', true);
let zoomLevel = store.get('zoomLevel', 100);

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
  console.log('DOM loaded, initializing app...');
  
  // Initialize highlight.js
  hljs.configure({
    languages: ['javascript', 'python', 'html', 'css', 'json', 'bash', 'sql', 'markdown', 'xml', 'yaml', 'go']
  });
  
  initializeTheme();
  initializeZoom();
  loadRecentFiles();
  setupEventListeners();
});

// Theme management
function initializeTheme() {
  setTheme(isDarkTheme);
}

function setTheme(dark) {
  isDarkTheme = dark;
  document.body.classList.toggle('dark-theme', dark);
  themeIcon.textContent = dark ? '☀️' : '🌙';
  
  // Update highlight.js theme
  const themePath = dark ? 
    '../../node_modules/highlight.js/styles/github-dark.css' : 
    '../../node_modules/highlight.js/styles/github.css';
  highlightTheme.href = themePath;
  
  // Store preference
  store.set('darkTheme', dark);
}

function toggleTheme() {
  setTheme(!isDarkTheme);
}

// Zoom management
function initializeZoom() {
  setZoom(zoomLevel);
}

function setZoom(level) {
  zoomLevel = Math.max(50, Math.min(200, level)); // Clamp between 50% and 200%
  
  if (markdownContent) {
    markdownContent.style.fontSize = `${zoomLevel}%`;
  }
  
  if (zoomLevelSpan) {
    zoomLevelSpan.textContent = `${zoomLevel}%`;
  }
  
  // Store preference
  store.set('zoomLevel', zoomLevel);
}

function zoomIn() {
  setZoom(zoomLevel + 10);
}

function zoomOut() {
  setZoom(zoomLevel - 10);
}

function resetZoom() {
  setZoom(100);
}

// File operations
async function openFile() {
  try {
    const result = await ipcRenderer.invoke('show-open-dialog');
    if (result.success && result.filePath) {
      loadFileFromPath(result.filePath);
    }
  } catch (error) {
    console.error('Error opening file:', error);
  }
}

function loadMarkdownFile(fileData) {
  const { content, filePath, fileName } = fileData;
  
  try {
    // Convert markdown to HTML
    const html = marked(content);
    
    // Update UI
    markdownContent.innerHTML = html;
    currentFileSpan.textContent = fileName;
    currentFilePath = filePath;
    
    // Apply syntax highlighting to all code blocks
    hljs.highlightAll();
    
    // Show content view
    welcomeScreen.classList.add('hidden');
    contentView.classList.remove('hidden');
    
    // Scroll to top when opening a new file
    markdownContent.scrollTop = 0;
    
    // Add to recent files
    addToRecentFiles(filePath, fileName);
    
  } catch (error) {
    console.error('Error rendering markdown:', error);
    markdownContent.innerHTML = `<p style="color: red;">Error rendering markdown: ${error.message}</p>`;
  }
}

// Recent files management
function addToRecentFiles(filePath, fileName) {
  let recentFiles = store.get('recentFiles', []);
  
  // Remove if already exists
  recentFiles = recentFiles.filter(file => file.path !== filePath);
  
  // Add to beginning
  recentFiles.unshift({
    path: filePath,
    name: fileName,
    timestamp: Date.now()
  });
  
  // Keep only last 10 files
  recentFiles = recentFiles.slice(0, 10);
  
  // Store and update UI
  store.set('recentFiles', recentFiles);
  loadRecentFiles();
}

function loadRecentFiles() {
  const recentFiles = store.get('recentFiles', []);
  
  if (recentFiles.length === 0) {
    recentFilesList.innerHTML = '<p class="no-recent">No recent files</p>';
    return;
  }
  
  recentFilesList.innerHTML = recentFiles.map(file => `
    <div class="recent-file-item" data-path="${file.path}">
      <div class="recent-file-name">${file.name}</div>
      <div class="recent-file-path">${file.path}</div>
    </div>
  `).join('');
  
  // Add click listeners to recent files
  document.querySelectorAll('.recent-file-item').forEach(item => {
    item.addEventListener('click', () => {
      const filePath = item.dataset.path;
      loadFileFromPath(filePath);
    });
  });
}

async function loadFileFromPath(filePath) {
  try {
    const result = await ipcRenderer.invoke('read-file', filePath);
    
    if (result.success) {
      loadMarkdownFile({
        content: result.content,
        filePath: filePath,
        fileName: filePath.split(/[\\/]/).pop()
      });
    } else {
      console.error('Error reading file:', result.error);
    }
  } catch (error) {
    console.error('Error loading file:', error);
  }
}

// Drag and drop functionality
function setupDragAndDrop() {
  let dragCounter = 0;
  
  document.addEventListener('dragenter', (e) => {
    e.preventDefault();
    dragCounter++;
    dropOverlay.classList.remove('hidden');
  });
  
  document.addEventListener('dragleave', (e) => {
    e.preventDefault();
    dragCounter--;
    if (dragCounter === 0) {
      dropOverlay.classList.add('hidden');
    }
  });
  
  document.addEventListener('dragover', (e) => {
    e.preventDefault();
  });
  
  document.addEventListener('drop', (e) => {
    e.preventDefault();
    dragCounter = 0;
    dropOverlay.classList.add('hidden');
    
    const files = Array.from(e.dataTransfer.files);
    const markdownFile = files.find(file => 
      file.name.endsWith('.md') || 
      file.name.endsWith('.markdown') || 
      file.name.endsWith('.txt')
    );
    
    if (markdownFile) {
      const reader = new FileReader();
      reader.onload = (event) => {
        loadMarkdownFile({
          content: event.target.result,
          filePath: markdownFile.path || markdownFile.name,
          fileName: markdownFile.name
        });
      };
      reader.readAsText(markdownFile);
    }
  });
}

// Event listeners
function setupEventListeners() {
  console.log('Setting up event listeners...');
  console.log('openFileBtn:', openFileBtn);
  console.log('welcomeOpenBtn:', welcomeOpenBtn);
  
  // Theme toggle
  if (themeToggle) {
    themeToggle.addEventListener('click', toggleTheme);
  }
  
  // File open buttons
  if (openFileBtn) {
    openFileBtn.addEventListener('click', () => {
      console.log('Open file button clicked');
      openFile();
    });
  }
  
  if (welcomeOpenBtn) {
    welcomeOpenBtn.addEventListener('click', () => {
      console.log('Welcome open button clicked');
      openFile();
    });
  }
  
  // Zoom controls
  if (zoomInBtn) {
    zoomInBtn.addEventListener('click', zoomIn);
  }
  
  if (zoomOutBtn) {
    zoomOutBtn.addEventListener('click', zoomOut);
  }
  
  if (zoomResetBtn) {
    zoomResetBtn.addEventListener('click', resetZoom);
  }
  
  // Drag and drop
  setupDragAndDrop();
  
  // Keyboard shortcuts
  document.addEventListener('keydown', (e) => {
    if (e.ctrlKey || e.metaKey) {
      switch (e.key) {
        case 'o':
          e.preventDefault();
          openFile();
          break;
        case 'd':
          e.preventDefault();
          toggleTheme();
          break;
        case '=':
        case '+':
          e.preventDefault();
          zoomIn();
          break;
        case '-':
          e.preventDefault();
          zoomOut();
          break;
        case '0':
          e.preventDefault();
          resetZoom();
          break;
      }
    }
  });
}

// IPC listeners
ipcRenderer.on('load-markdown', (event, fileData) => {
  loadMarkdownFile(fileData);
});

ipcRenderer.on('toggle-theme', () => {
  toggleTheme();
});

// Export functions for potential use
window.electronAPI = {
  loadMarkdownFile,
  toggleTheme,
  openFile
};