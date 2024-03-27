function toggleSidebar() {
    const sidebar = document.querySelector('.sidebar');
    sidebar.classList.toggle('active');
}

// Add an event listener to close the sidebar when clicking outside of it
document.addEventListener('click', function (event) {
    const sidebar = document.querySelector('.sidebar');
    if (sidebar.classList.contains('active') && !event.target.closest('.sidebar') && !event.target.closest('.nav-icon')) {
        sidebar.classList.remove('active');
    }
});

const selectImage = document.querySelector('.select-image');
const inputFile = document.querySelector('#file');
const imgArea = document.querySelector('.img-area');

selectImage.addEventListener('click', function () {
    inputFile.click();
})

const section = document.querySelector("section"),
    overlay = document.querySelector(".overlay"),
    showBtn = document.querySelector(".buttonSubmit"),
    closeBtn = document.querySelector(".close-btn");

showBtn.addEventListener("click", () => {
    // Validate the form fields
    const isAnonymousChecked = document.getElementById('is-anonymous').checked;
    const typeOfIncident = document.getElementById('type-of-incident').value.trim();
    const location = document.getElementById('location').value.trim();
    const useLocationChecked = document.getElementById('useLocation').checked;
    const description = document.getElementById('description').value.trim();

    if (!isAnonymousChecked && typeOfIncident === '' && location === '' && !useLocationChecked && description === '') {
        alert('Please fill in all required fields.');
        return;
    }

    const isConfirmed = window.confirm("Are you sure you want to submit the report?");

if (isConfirmed) {
    // Clear input fields
    document.getElementById('is-anonymous').checked = false;
    document.getElementById('type-of-incident').value = '';
    document.getElementById('location').value = '';
    document.getElementById('useLocation').checked = false;
    document.getElementById('description').value = '';

    // Clear image area
    imgArea.innerHTML = '<i class=\'bx bxs-cloud-upload icon\'></i><h3>Upload Image</h3><p>Image size must be less than <span>5MB</span></p>';
    imgArea.classList.remove('active');
    imgArea.dataset.img = '';

    // Show the success modal
    section.classList.add("active");
    document.querySelector('.modal-box').setAttribute('aria-hidden', 'false');
    document.querySelector('.overlay').setAttribute('aria-hidden', 'false');
    section.classList.add("active");
    document.querySelector('.modal-box').setAttribute('aria-hidden', 'false');
    document.querySelector('.overlay').setAttribute('aria-hidden', 'false');
}
});

overlay.addEventListener("click", () =>
    section.classList.remove("active")
);

closeBtn.addEventListener("click", () =>
    section.classList.remove("active")
);

inputFile.addEventListener('change', function () {
    const image = this.files[0]
    if (image.size < 5000000) {
        const reader = new FileReader();
        reader.onload = () => {
            const allImg = imgArea.querySelectorAll('img');
            allImg.forEach(item => item.remove());
            const imgUrl = reader.result;
            const img = document.createElement('img');
            img.src = imgUrl;
            imgArea.appendChild(img);
            imgArea.classList.add('active');
            imgArea.dataset.img = image.name;
        }
        reader.readAsDataURL(image);
    } else {
        alert("Image size more than 5MB");
    }
})


function handleIncidentType() {
    var incidentType = document.getElementById("incident_type").value;
    var otherIncidentTypeContainer = document.getElementById("other_incident_type_container");

    if (incidentType === "Other") {
        otherIncidentTypeContainer.style.display = "block";
        otherIncidentTypeContainer.getElementsByTagName("input")[0].setAttribute("name", "type_of_incident");
    } else {
        otherIncidentTypeContainer.style.display = "none";
        otherIncidentTypeContainer.getElementsByTagName("input")[0].removeAttribute("name");
    }
}

function submitForm(event) {
    event.preventDefault(); // Prevent the default form submission

    // Your form submission logic goes here
    // For demonstration, we'll just show a toast notification
    showToast("Report successfully submitted");
}

function showToast(message) {
    // Create a new div element
    const toast = document.createElement("div");
    toast.className = "toast";
    toast.textContent = message;

    // Append the toast to the body
    document.body.appendChild(toast);

    // Remove the toast after a certain duration (e.g., 3 seconds)
    setTimeout(function() {
        toast.remove();
    }, 3000);
}