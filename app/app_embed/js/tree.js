document.addEventListener("DOMContentLoaded", () => {
    document.querySelectorAll("#tree .node").forEach(node => {
        const children = node.querySelector(".children");
        if (children) {
            node.addEventListener("click", (e) => {
                e.stopPropagation();
                node.classList.toggle("expanded");
            });
        }
    });
});
