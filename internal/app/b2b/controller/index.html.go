package controller

// This file is generated from the front end project b2b-client.
// Do not touch it.

// B2BIndexHTML is the index file for angular web app.
const B2BIndexHTML = `
<!DOCTYPE html>
<html lang="en">

<head>
    <base href="/corporate">
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <link rel="apple-touch-icon" sizes="180x180" href="http://interactive.ftchinese.com/favicons/apple-touch-icon-180x180.png">
    <link rel="apple-touch-icon" sizes="152x152" href="http://interactive.ftchinese.com/favicons/apple-touch-icon-152x152.png">
    <link rel="apple-touch-icon" sizes="120x120" href="http://interactive.ftchinese.com/favicons/apple-touch-icon-120x120.png">
    <link rel="apple-touch-icon" sizes="76x76" href="http://interactive.ftchinese.com/favicons/apple-touch-icon-76x76.png">
    <link href="http://interactive.ftchinese.com/favicons/favicon.ico" type="image/x-icon" rel="shortcut icon" />
    <title>B2B订阅 - FT中文网</title>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.5.0/css/bootstrap.min.css" rel="stylesheet">
    
    <style>app-toolbar,body{background-color:#fff1e5}body{font-family:Helvetica Neue,Lucida Grande,Verdana,PingFang SC,PingFang TC,Hiragino Sans GB,Heiti SC,Heiti TC,WenQuanYi Micro Hei,Microsoft YaHei,Microsoft JhengHei,STHeiti,sans-serif}.page-content{min-height:80vh}table caption{caption-side:top;text-align:center}a{color:#0d7680}a:hover,app-toolbar a:hover{color:#4d4845}.alert{border-radius:0}@media (min-width:1200px){.w-xl-50,.w-xl-75{width:50%;margin:auto}.w-xl-75{width:75%}}@media (min-width:992px){.w-lg-50,.w-lg-75{width:50%;margin:0 auto}.w-lg-75{width:75%}}app-toolbar{width:100%;height:68px;padding:16px;display:flex;flex-direction:row;justify-content:space-between;align-items:center;border-bottom:1px solid #ccc1b7}app-toolbar a{color:#33302e}app-toolbar .ftc-brand a:hover{text-decoration:none}app-toolbar .cart-button{display:inline-flex;align-items:center;padding:8px 16px;border-radius:2px;font-size:14px;margin-right:2em;cursor:pointer;background-color:#fff;position:relative}.badge-container .badge,app-toolbar .cart-button .badge{position:absolute;right:0;top:50%;transform:translate(50%,-50%)}app-toolbar .cart-button:hover{opacity:.8;text-decoration:none}.badge-container{position:relative}.badge-container .badge{top:0}.nav.sidebar a.active{color:#990f3d;border-right:2px solid #990f3d}.o-card-group .card{height:100%;background-color:transparent;border-radius:0}.o-card-group .card ul{padding-left:0}.o-card-group .card li{list-style-position:outside;list-style-type:none;padding-left:1.5em;margin-top:.5em;margin-bottom:.5em;background-size:1.5em;background-image:url(https://www.ft.com/__origami/service/image/v2/images/raw/fticon-v1:tick?format=svg&source=ftchinese&tint=26747a);background-repeat:no-repeat;background-position:0 0}.o-card-group .card-title{border-bottom:1px solid #26747a}.o-card-group .card-footer{background-color:transparent;border-top:0}.btn{border-radius:0}.btn-primary{background-color:#0d7680;color:#fff}.btn-primary.focus,.btn-primary:focus,.btn-primary:hover{background-color:#0a5e66;border-color:transparent}.btn-primary.disabled,.btn-primary:disabled{background-color:#0d7680;pointer-events:none;opacity:.4;cursor:default}.btn-primary:not(:disabled):not(.disabled).active,.btn-primary:not(:disabled):not(.disabled):active,.show>.btn-primary.dropdown-toggle{color:#fff;background-color:#40858b;border-color:transparent}.btn-link{color:#0d7680}.btn-outline-primary{color:#0d7680;background-color:transparent;border-color:#0d7680}.btn-outline-primary.focus,.btn-outline-primary:focus,.btn-outline-primary:hover{color:#0d7680;background-color:rgba(13,118,128,.1);border-color:#0d7680}.btn-outline-primary.disabled,.btn-outline-primary:disabled{pointer-events:none;opacity:.4;cursor:default}.btn-outline-primary:not(:disabled):not(.disabled).active,.btn-outline-primary:not(:disabled):not(.disabled):active,.show>.btn-outline-primary.dropdown-toggle{color:#fff;background-color:#0d7680;border-color:#0d7680}legend{border-bottom:1px solid #ccc1b7}.form-label{font-weight:600;color:#33302e}.form-control{border-radius:0}.form-control:disabled{background-color:#e6d9ce;color:#66605c;border-color:#e6d9ce;cursor:default}.form-errortext{color:#c00;font-size:14px;line-height:16px;clear:both;display:block;margin-top:-1px;padding:3px 0}.form-invalid{background-color:#fff;color:#c00;border-color:#c00}.o-forms-input__error{font-size:14px;line-height:16px;color:#c00;display:block;position:relative;margin-top:4px}.list-group-item{background-color:transparent}.barrier-banner{background-color:#f2e0d0}.barrier-banner .container{padding-bottom:20px;padding-top:20px}.barrier-cover,.box-16-9 img{width:100%;-o-object-fit:contain;object-fit:contain}.barrier-heading{font-size:32px;line-height:32px;font-weight:600;opacity:.5}.box-16-9{position:relative;display:block;height:0;overflow:hidden;padding:0 0 56.25%;background-color:#f2dfce}.box-16-9 img{position:absolute;top:0;left:0;height:100%}.o-footer{margin-top:40px;border-top:10px solid;line-height:20px;font-size:16px}.o-footer a{border-bottom:0;text-decoration:none}.o-footer__matrix-title{margin:0;line-height:inherit;font-size:inherit;font-weight:600}.o-footer__matrix-content{margin-top:10px;margin-bottom:10px}.o-footer__matrix-link{display:block;padding:4px 0}.o-footer--theme-dark{color:#fff;background-color:#33302e;border-color:#0a5e66}.o-footer--theme-dark .o-footer__copyright,.o-footer--theme-dark a{color:#b3a9a0}.o-footer--theme-dark a:focus,.o-footer--theme-dark a:hover,.o-footer--theme-light a:focus,.o-footer--theme-light a:hover{color:#fff}.o-footer__link-items{padding:20px 0;list-style-type:none;margin:0}.o-footer__link-item{display:inline-block;margin-bottom:10px;margin-right:20px}.o-footer--theme-light{color:#33302e;background-color:#f2e5da;border-color:#990f3d}.o-footer--theme-light .o-footer__copyright,.o-footer--theme-light a{color:#4d4845}@media print{.o-footer{display:none}}.modal{display:block;padding-right:15px;overflow-x:hidden;overflow-y:auto;z-index:1072}</style>
    
</head>

<body>
    <app-root></app-root>
    
    <script src="/frontend/b2b/runtime-es2015.409e6590615fb48d139f.js" type="module"></script>
    
    <script src="/frontend/b2b/runtime-es5.409e6590615fb48d139f.js" nomodule defer></script>
    
    <script src="/frontend/b2b/polyfills-es5.341d48ffc88eaeae911d.js" nomodule defer></script>
    
    <script src="/frontend/b2b/polyfills-es2015.95bb65e853d2781ac3bd.js" type="module"></script>
    
    <script src="/frontend/b2b/main-es2015.bb0ae2c9057ca1cc2111.js" type="module"></script>
    
    <script src="/frontend/b2b/main-es5.bb0ae2c9057ca1cc2111.js" nomodule defer></script>
    
</body>
</html>
`
