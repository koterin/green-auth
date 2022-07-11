'use stirct';

const switcher = document.getElementById('tglBtn');

document.onload = () => {
    let theme = localStorage.getItem('theme');
  //  if (theme != null) {
        document.body.className = theme;
//    }
}

// Switch between dark and light modes
switcher.addEventListener('click', function() {
    document.body.classList.toggle('dark-theme');
    document.body.classList.toggle('light-theme');

    const className = document.body.className;

    if (className == "light-theme") {
        this.textContent = "Dark";
    } else {
        this.textContent = "Light";
    }

    localStorage.setItem('theme', className);
});
