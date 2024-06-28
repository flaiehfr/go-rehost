# go-rehost

Utilitaire en ligne de commande pour uploader une image depuis votre ordinateur sur Rehost.

Prends un paramètre `-file` désignant le fichier à uploader, et écrit sur la sortie standard l'URL du fichier.

### Exemple en anonyme:
```
> .\go-rehost.exe -file C:/temp/forum_logo.gif
https://rehost.diberie.com/Picture/Get/f/296858
```

### Exemple avec votre cookie:

Ouvrez votre browser favori et récupérez la valeur du cookie `.AspNet.ApplicationCookie` pour `rehost.diberie.com` (`F12` > `Application` > `Cookies` > `https://rehost.diberie.com`)

![https://rehost.diberie.com/Picture/Get/f/296980](https://rehost.diberie.com/Picture/Get/f/296980)

Passez le à l'aide du paramètre `-cookie`

```
> .\go-rehost.exe -file C:/temp/forum_logo.gif -cookie vyqYB4ldt0plNN5P17sU7k7Vr9ROr-0XLAPo065C2d602XK_N5pCT-7Q3Yp...
https://rehost.diberie.com/Picture/Get/f/296978
```
