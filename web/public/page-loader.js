// page-loader.js — skeleton fullscreen loader for Uszata Kuchnia SPA
// Dismiss: window.hideLoader() after app mounts

(function () {
    const dark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const bg = dark ? '#0A0A0F' : '#F8F8FF';
    const bgAlt = dark ? '#13131C' : '#FFFFFF';
    const rule = dark ? 'rgba(248,248,255,0.08)' : 'rgba(10,10,15,0.08)';
    const s1 = dark ? '#1E1E2A' : '#EDEDF5';
    const s2 = dark ? '#26263A' : '#F4F4FA';
    const ink = dark ? '#F8F8FF' : '#0A0A0F';
    const acc = dark ? '#9D8FE8' : '#6B4EE6';

    const path = window.location.pathname;
    const isRecipe =
        /^\/recipes\/\d+$/.test(path) ||
        /^\/recipes\/\d+\/edit$/.test(path);
    const isNew = path === '/recipes/new';
    const isPanel = path.startsWith('/admin');

    const css = `
    #__pl {
      position: fixed; inset: 0; z-index: 9999;
      background: ${bg};
      transition: opacity .3s ease;
      font-family: "Geist","Inter Tight",system-ui,sans-serif;
    }
    #__pl .topbar {
      height: 52px;
      background: ${bgAlt};
      border-bottom: 1px solid ${rule};
      display: flex; align-items: center;
      padding: 0 32px; gap: 32px;
    }
    #__pl .logo {
      font-size: 17px; font-weight: 600; letter-spacing: -0.4px; color: ${ink};
    }
    #__pl .logo span { color: ${acc}; }
    #__pl .nav { display: flex; gap: 20px; }
    #__pl .nav-item {
      width: 48px; height: 10px; border-radius: 4px;
    }
    #__pl .body { padding: 40px 56px; }
    #__pl .body-panel { display: grid; grid-template-columns: 200px 1fr; height: calc(100vh - 52px); }
    #__pl .sidebar { background: ${bgAlt}; border-right: 1px solid ${rule}; padding: 16px 12px; display: flex; flex-direction: column; gap: 8px; }
    #__pl .sidebar-item { height: 32px; border-radius: 7px; }
    #__pl .main { padding: 28px 32px; display: flex; flex-direction: column; gap: 14px; }

    @keyframes __sh { 0%{background-position:-600px 0} 100%{background-position:600px 0} }
    #__pl .sk {
      background: linear-gradient(90deg, ${s1} 25%, ${s2} 50%, ${s1} 75%);
      background-size: 1200px 100%;
      animation: __sh 1.6s ease-in-out infinite;
      border-radius: 6px;
    }

    #__pl .detail { display: grid; grid-template-columns: 1fr 420px; gap: 48px; padding: 40px 56px; }
    #__pl .detail-left { display: flex; flex-direction: column; gap: 14px; }
    #__pl .detail-img { border-radius: 10px; }
  `;

    function listSkeleton() {
        let cards = '';
        for (let i = 0; i < 6; i++) {
            const w = [55, 70, 60, 80, 50, 65][i];
            cards += `
        <div style="display:flex;gap:16px;padding:0 0 18px;border-bottom:1px solid ${rule};margin-bottom:18px;">
          <div class="sk" style="width:100px;height:100px;border-radius:8px;flex-shrink:0;"></div>
          <div style="flex:1;display:flex;flex-direction:column;gap:8px;padding-top:4px;">
            <div class="sk" style="height:16px;width:${w}%;border-radius:5px;"></div>
            <div class="sk" style="height:11px;width:40%;border-radius:4px;"></div>
            <div class="sk" style="height:11px;width:55%;border-radius:4px;margin-top:4px;"></div>
          </div>
        </div>`;
        }
        return `
      <div class="body">
        <div style="display:grid;grid-template-columns:260px 1fr;gap:40px;">
          <div style="display:flex;flex-direction:column;gap:10px;">
            <div class="sk" style="height:14px;width:70%;border-radius:5px;margin-bottom:6px;"></div>
            ${[80, 65, 75, 55, 70].map((w) => `<div class="sk" style="height:30px;border-radius:7px;width:${w}%;"></div>`).join('')}
          </div>
          <div>${cards}</div>
        </div>
      </div>`;
    }

    function recipeSkeleton() {
        return `
      <div class="detail">
        <div class="detail-left">
          <div class="sk" style="height:13px;width:30%;border-radius:4px;"></div>
          <div class="sk" style="height:36px;width:80%;border-radius:6px;margin-top:4px;"></div>
          <div class="sk" style="height:13px;width:60%;border-radius:4px;"></div>
          <div class="sk" style="height:48px;border-radius:8px;margin-top:8px;"></div>
          <div style="display:flex;flex-direction:column;gap:10px;margin-top:16px;">
            ${[100, 90, 95, 80].map((w) => `<div class="sk" style="height:13px;width:${w}%;border-radius:4px;"></div>`).join('')}
          </div>
        </div>
        <div class="sk detail-img" style="height:520px;"></div>
      </div>`;
    }

    function newSkeleton() {
        return `
      <div class="body" style="max-width:760px;margin:0 auto;">
        <div class="sk" style="height:28px;width:40%;border-radius:6px;margin-bottom:28px;"></div>
        ${[1, 2, 3]
            .map(
                () => `
          <div style="background:${bgAlt};border:1px solid ${rule};border-radius:10px;padding:20px 24px;margin-bottom:16px;">
            <div class="sk" style="height:14px;width:25%;border-radius:4px;margin-bottom:16px;"></div>
            <div class="sk" style="height:40px;border-radius:7px;"></div>
          </div>`,
            )
            .join('')}
      </div>`;
    }

    function panelSkeleton() {
        return `
      <div class="body-panel">
        <div class="sidebar">
          ${[70, 55, 65, 45, 50].map((w) => `<div class="sk sidebar-item" style="width:${w}%;"></div>`).join('')}
        </div>
        <div class="main">
          <div class="sk" style="height:20px;width:20%;border-radius:5px;"></div>
          <div style="background:${bgAlt};border:1px solid ${rule};border-radius:9px;overflow:hidden;">
            ${[100, 85, 90, 75, 95, 80]
                .map(
                    (w) => `
              <div style="display:flex;gap:16px;padding:10px 16px;border-bottom:1px solid ${rule};">
                <div class="sk" style="flex:1;height:13px;width:${w}%;border-radius:4px;"></div>
                <div class="sk" style="width:60px;height:13px;border-radius:4px;"></div>
                <div class="sk" style="width:50px;height:13px;border-radius:4px;"></div>
              </div>`,
                )
                .join('')}
          </div>
        </div>
      </div>`;
    }

    const topbar = `
    <div class="topbar">
      <div class="logo">uszatakuchnia<span>.</span></div>
      <div class="nav">
        <div class="sk nav-item"></div>
        <div class="sk nav-item"></div>
        <div class="sk nav-item"></div>
      </div>
    </div>`;

    let body = '';
    if (isRecipe) body = recipeSkeleton();
    else if (isNew) body = newSkeleton();
    else if (isPanel) body = panelSkeleton();
    else body = listSkeleton();

    const style = document.createElement('style');
    style.textContent = css;
    document.head.appendChild(style);

    const el = document.createElement('div');
    el.id = '__pl';
    el.innerHTML = topbar + body;
    document.body.appendChild(el);

    window.hideLoader = function () {
        const l = document.getElementById('__pl');
        if (!l) return;
        l.style.opacity = '0';
        setTimeout(() => l.remove(), 320);
    };
})();
