---
title: "Avatar"
imports:
  - rh-avatar
---

<p>9 demos for <code>&lt;rh-avatar&gt;</code></p>


### color context

<rh-context-demo>
  <rh-avatar>George Boole
    <span slot="subtitle">Professor of Mathematics, </span>
    <a slot="subtitle"
       href="https://www.wikiwand.com/en/Queen's_College,_Cork">Queen's College, Cork</a>
  </rh-avatar>
  <rh-avatar name="John von Neumann"
             subtitle="Mathematician"
             plain></rh-avatar>
  <rh-avatar name="Grace Hopper"
             subtitle="Rear Admiral"
             src="https://ux.redhat.com/elements/avatar/demo/hopper.jpg"
             plain></rh-avatar>
  <rh-avatar name="Haskell Curry"
             subtitle="Computer Scientist"
             pattern="squares"
             plain></rh-avatar>
  <rh-avatar name="Edsger Dijkstra"
             subtitle="Computer Scientist"
             pattern="triangles"
             plain></rh-avatar>
</rh-context-demo>


### index

<rh-avatar name="Omar Khayyam"
           subtitle="Mathematician, Astronomer"
           src="https://upload.wikimedia.org/wikipedia/commons/thumb/0/08/005-a-Ruby-kindles-in-the-vine-810x1146.jpg/212px-005-a-Ruby-kindles-in-the-vine-810x1146.jpg"></rh-avatar>


### links

<figure>
  <figcaption>Links applied to Name</figcaption>
  <rh-avatar src="https://upload.wikimedia.org/wikipedia/commons/thumb/7/77/Jeannette_Wing%2C_Davos_2013.jpg/330px-Jeannette_Wing%2C_Davos_2013.jpg">
    <a href="https://www.wikiwand.com/en/Jeannette_Wing">Jeannette Wing</a>
    <span slot="subtitle">Avanessians Director of the Data Science Institute, Columbia University</span>
  </rh-avatar>
</figure>
<figure>
  <figcaption>Links applied to job details</figcaption>
  <rh-avatar src="https://upload.wikimedia.org/wikipedia/commons/thumb/c/ce/George_Boole_color.jpg/330px-George_Boole_color.jpg">George Boole
    <span slot="subtitle">Professor of Mathematics, </span>
    <a slot="subtitle" href="https://www.wikiwand.com/en/Queen's_College,_Cork">Queen's College, Cork</a>
  </rh-avatar>
</figure>


### pattern

<figure>
  <figcaption>Squares pattern</figcaption>
  <rh-avatar name="Alonzo Church"
             subtitle="Inventor of the Lambda Calculus"
             pattern="squares"></rh-avatar>
</figure>
<figure>
  <figcaption>Triangles pattern</figcaption>
  <rh-avatar name="Alan Turing"
             subtitle="Cryptographer"
             pattern="triangles"></rh-avatar>
</figure>


### plain

<section id="plain-avatar">
  <rh-avatar name="John von Neumann"
             subtitle="Mathematician"
             plain></rh-avatar>
  <rh-avatar name="Grace Hopper"
             subtitle="Rear Admiral"
             src="https://ux.redhat.com/elements/avatar/demo/hopper.jpg"
             plain></rh-avatar>
  <rh-avatar name="Haskell Curry"
             subtitle="Computer Scientist"
             pattern="squares"
             plain></rh-avatar>
  <rh-avatar name="Edsger Dijkstra"
             subtitle="Computer Scientist"
             pattern="triangles"
             plain></rh-avatar>
</section>


### position

<rh-avatar name="Ada Lovelace"
           subtitle="Computer Programmer"
           layout="block"></rh-avatar>


### sizes

<figure>
  <figcaption><code>--rh-size-icon-03</code></figcaption>
  <rh-avatar name="Radia Perlman"
             subtitle="Mother of the Internet"
             src="https://upload.wikimedia.org/wikipedia/commons/thumb/a/af/Radia_Perlman_2009.jpg/330px-Radia_Perlman_2009.jpg"
             style="--rh-avatar-size:var(--rh-size-icon-03, 32px);">
  </rh-avatar>
</figure>
<figure>
  <figcaption><code>--rh-size-icon-05</code></figcaption>
  <rh-avatar src="https://upload.wikimedia.org/wikipedia/commons/thumb/e/eb/Gordon_Moore_1978_%28cropped%29.png/330px-Gordon_Moore_1978_%28cropped%29.png"
             name="Gordon Moore"
             style="--rh-avatar-size:var(--rh-size-icon-05, 48px);">
    <span slot="subtitle">Co-founder, <em>Intel</em></span>
  </rh-avatar>
</figure>
<figure>
  <figcaption><code>--rh-size-icon-06</code> <small>(default)</small></figcaption>
  <rh-avatar src="https://upload.wikimedia.org/wikipedia/commons/thumb/6/6d/Katherine_Johnson_1983.jpg/330px-Katherine_Johnson_1983.jpg"
             name="Katherine Johnson"
             subtitle="Recipient, National Medal of Freedom 2016">
  </rh-avatar>
</figure>
<figure>
  <figcaption>
    <code>--rh-size-icon-08</code>
    <small>Avatars cannot be larger than <code>--rh-size-icon-06</code></small>
  </figcaption>
  <rh-avatar name="Hedy Lamarr"
             src="https://upload.wikimedia.org/wikipedia/commons/thumb/8/83/Hedy_Lamarr_Publicity_Photo_for_The_Heavenly_Body_1944.jpg/330px-Hedy_Lamarr_Publicity_Photo_for_The_Heavenly_Body_1944.jpg"
             subtitle="Jewish actress and inventor"
             style="--rh-avatar-size:var(--rh-size-icon-08, 96px);"></rh-avatar>
</figure>


### subtitles

<rh-avatar src="https://ux.redhat.com/elements/avatar/demo/schoenfinkel.jpg">Moses Schoenfinkle
  <span slot="subtitle">
    Inventor of Combinatorics,
    often uncreditted for inventing the process of "currying" functions,
    however, "schoenfinkling" doesn't exactly roll off the tongue, so we'll
    let it slide
  </span>
</rh-avatar>


### variants

<section id="bordered-avatars">
  <rh-avatar name="Bordered von Neumann"
             subtitle="Mathematician"
             variant="bordered"></rh-avatar>
  <rh-avatar name="Grace Bordered Jr."
             subtitle="Rear Admiral"
             src="https://ux.redhat.com/elements/avatar/demo/hopper.jpg"
             variant="bordered"></rh-avatar>
</section>

