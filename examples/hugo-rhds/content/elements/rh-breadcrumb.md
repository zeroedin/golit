---
title: "Breadcrumb"
imports:
  - rh-breadcrumb
lightdom:
  - rh-breadcrumb-lightdom.css
---

<p>6 demos for <code>&lt;rh-breadcrumb&gt;</code></p>


### color context

<rh-context-demo>
  <rh-breadcrumb>
    <ol>
      <li><a href="#">Home</a></li>
      <li><a href="#">Products</a></li>
      <li><a href="#">Red Hat OpenShift on AWS</a></li>
      <li><a href="#">4</a></li>
      <li><a href="#">Introduction to ROSA</a></li>
      <li><a href="#" aria-current="page">Chapter 1. Understanding ROSA</a></li>
    </ol>
  </rh-breadcrumb>
  <h2 class="subtle-heading">Subtle:</h2>
  <rh-breadcrumb variant="subtle">
    <ol>
      <li><a href="#">Home</a></li>
      <li><a href="#">Products</a></li>
      <li><a href="#">Red Hat OpenShift on AWS</a></li>
      <li><a href="#">4</a></li>
      <li><a href="#">Introduction to ROSA</a></li>
      <li><a href="#" aria-current="page">Chapter 1. Understanding ROSA</a></li>
    </ol>
  </rh-breadcrumb>
</rh-context-demo>


### custom accessible label

<rh-breadcrumb accessible-label="Breadcrumb Navigation">
  <ol>
    <li><a href="#">Home</a></li>
    <li><a href="#">Products</a></li>
    <li><a href="#">Red Hat OpenShift on AWS</a></li>
    <li><a href="#">4</a></li>
    <li><a href="#">Introduction to ROSA</a></li>
    <li><a href="#" aria-current="page">Chapter 1. Understanding ROSA</a></li>
  </ol>
</rh-breadcrumb>


### index

<rh-breadcrumb>
  <ol>
    <li><a href="#">Home</a></li>
    <li><a href="#">Products</a></li>
    <li><a href="#">Red Hat OpenShift on AWS</a></li>
    <li><a href="#">4</a></li>
    <li><a href="#">Introduction to ROSA</a></li>
    <li><a href="#" aria-current="page">Chapter 1. Understanding ROSA</a></li>
  </ol>
</rh-breadcrumb>


### non interactive last item

<rh-breadcrumb>
  <ol>
    <li><a href="#">Home</a></li>
    <li><a href="#">Products</a></li>
    <li><a href="#">Red Hat OpenShift on AWS</a></li>
    <li><a href="#">4</a></li>
    <li><a href="#">Introduction to ROSA</a></li>
    <li>Chapter 1. Understanding ROSA</li>
  </ol>
</rh-breadcrumb>
<h2 class="subtle-heading">Subtle:</h2>
<rh-breadcrumb variant="subtle">
  <ol>
    <li><a href="#">Home</a></li>
    <li><a href="#">Products</a></li>
    <li><a href="#">Red Hat OpenShift on AWS</a></li>
    <li><a href="#">4</a></li>
    <li><a href="#">Introduction to ROSA</a></li>
    <li>Chapter 1. Understanding ROSA</li>
  </ol>
</rh-breadcrumb>


### subtle

<rh-breadcrumb variant="subtle">
  <ol>
    <li><a href="#">Home</a></li>
    <li><a href="#">Products</a></li>
    <li><a href="#">Red Hat OpenShift on AWS</a></li>
    <li><a href="#">4</a></li>
    <li><a href="#">Introduction to ROSA</a></li>
    <li><a href="#" aria-current="page">Chapter 1. Understanding ROSA</a></li>
  </ol>
</rh-breadcrumb>


### truncate

<rh-breadcrumb truncate>
  <ol>
    <li><a href="#home">Home</a></li>
    <li><a href="#products">Products</a></li>
    <li><a href="#open-shift-aws">Red Hat OpenShift on AWS</a></li>
    <li><a href="#4">4</a></li>
    <li><a href="#introduction-to-rosa">Introduction to ROSA</a></li>
    <li><a href="#understanding-rosa" aria-current="page">Chapter 1. Understanding ROSA</a></li>
  </ol>
</rh-breadcrumb>

