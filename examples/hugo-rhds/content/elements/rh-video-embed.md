---
title: "Video embed"
imports:
  - rh-card
  - rh-cta
  - rh-video-embed
lightdom:
  - rh-cta-lightdom-shim.css
---

<p>6 demos for <code>&lt;rh-video-embed&gt;</code></p>


### alignment

<div class="wrap">
  <rh-video-embed class="centered">
    <img slot="thumbnail" src="../video-thumb.jpg" alt="Image description"/>
    <template>
      <iframe title="Title of video" width="900" height="499" src="https://www.youtube.com/embed/Hc8emNr2igU" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
    </template>
    <p slot="caption"><a class="rh-video-embed-caption-link" href="https://www.redhat.com/">View the infographic</a></p>
  </rh-video-embed>
</div>
<div class="wrap">
  <rh-video-embed class="right-aligned">
    <img slot="thumbnail" src="../video-thumb.jpg" alt="Image description"/>
    <template>
      <iframe title="Title of video" width="900" height="499" src="https://www.youtube.com/embed/Hc8emNr2igU" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
    </template>
    <p slot="caption"><a class="rh-video-embed-caption-link" href="https://www.redhat.com/">View the infographic</a></p>
  </rh-video-embed>
</div>


### card with video

<div class="wrapper">
  <rh-card>
    <rh-video-embed slot="header">
      <img slot="thumbnail" src="../video-thumb.jpg" alt="Image description"/>
      <template>
        <iframe title="Title of video" width="900" height="499" src="https://www.youtube.com/embed/Hc8emNr2igU" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
      </template>
    </rh-video-embed>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      Nullam eleifend elit sed est egestas, a sollicitudin mauris
      tincidunt. Pellentesque vel dapibus risus. Nullam aliquam
      felis orci, eget cursus mi lacinia quis. Vivamus at felis sem.</p>
    <rh-cta variant="secondary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
</div>


### color context

<rh-context-demo>
  <rh-video-embed>
    <img slot="thumbnail" src="../video-thumb.jpg" alt="Image description"/>
    <template>
      <iframe title="Title of video" width="900" height="499" src="https://www.youtube.com/embed/Hc8emNr2igU" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
    </template>
    <p slot="caption">Video caption here</p>
  </rh-video-embed>
</rh-context-demo>


### index

<rh-video-embed>
  <img slot="thumbnail" src="video-thumb.jpg" alt="Image description"/>
  <template>
    <iframe title="Title of video" width="900" height="499" src="https://www.youtube.com/embed/Hc8emNr2igU" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
  </template>
  <p slot="caption"><a class="rh-video-embed-caption-link" href="https://www.redhat.com/">View the infographic</a></p>
</rh-video-embed>


### no caption

<rh-video-embed>
  <img slot="thumbnail" src="../video-thumb.jpg" alt="Image description"/>
  <template>
    <iframe title="Title of video" width="900" height="499" src="https://www.youtube.com/embed/Hc8emNr2igU" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
  </template>
</rh-video-embed>


### require consent

<rh-video-embed id="video" require-consent>
  <img slot="thumbnail" src="../video-thumb.jpg" alt="Image description"/>
  <template>
    <iframe title="Title of video" width="900" height="499" src="https://www.youtube.com/embed/Hc8emNr2igU" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
  </template>
  <p slot="caption"><a class="rh-video-embed-caption-link" href="https://www.redhat.com/">View the infographic</a></p>
</rh-video-embed>

