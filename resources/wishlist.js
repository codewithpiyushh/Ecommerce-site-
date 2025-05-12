var wisharr = JSON.parse(localStorage.getItem("wishListObj")) || [];

var itemcount = wisharr.length;
document.querySelector(".wishcount").innerText = ` ${itemcount} Items`;

wisharr.map(function (ele, ind) {
    var box = document.createElement("div");
    box.className = "wishlist-item";

    // Image
    var image = document.createElement("img");
    image.src = ele.image_url;

    var imgbox = document.createElement("div");
    imgbox.className = "imgbox";
    imgbox.appendChild(image);

    // Brand name
    var para = document.createElement("p");
    para.innerText = ele.brand;
    para.style.color = "gray";

    // Price section
    var price = document.createElement("span");
    price.innerText = "Rs. " + ele.price;
    price.style.color = "black";

    var strikedprice = document.createElement("span");
    strikedprice.innerText = "Rs. " + ele.strikedoffprice;
    strikedprice.style.textDecoration = "line-through";
    strikedprice.style.color = "gray";

    var offer = document.createElement("span");
    offer.innerText = ele.offer;
    offer.style.color = "red";

    var pricepara = document.createElement("p");
    pricepara.className = "pricepara";
    pricepara.append(price, strikedprice, offer);

    // Buttons
    var buttonrm = document.createElement("button");
    buttonrm.innerText = "Remove";
    buttonrm.addEventListener("click", function () {
        removefromwish(ind);
    });

    var buttonbag = document.createElement("button");
    buttonbag.innerText = "MOVE TO BAG";
    buttonbag.addEventListener("click", function () {
        sendtocart(ele, ind);
    });

    var buttonholder = document.createElement("div");
    buttonholder.className = "buttonholder";
    buttonholder.append(buttonrm, buttonbag);

    // Content Box
    var contentBox = document.createElement("div");
    contentBox.className = "contentBox";
    contentBox.append(para, pricepara, buttonholder);

    // Final box append
    box.append(imgbox, contentBox);
    document.querySelector(".container").append(box);
});

function removefromwish(ind) {
    wisharr.splice(ind, 1);
    localStorage.setItem("wishListObj", JSON.stringify(wisharr));
    location.reload();  // Add this line to reflect the change
}

function sendtocart(ele, ind) {
    baglist.unshift(ele);
    localStorage.setItem("BagListObj", JSON.stringify(baglist));
    wisharr.splice(ind, 1);
    localStorage.setItem("wishListObj", JSON.stringify(wisharr));
    location.reload();  // Add this line too
}


var baglist = JSON.parse(localStorage.getItem("BagListObj")) || [];



document.getElementById("landingPage").addEventListener("click", function () {
    window.location.href = "index.html";
});
