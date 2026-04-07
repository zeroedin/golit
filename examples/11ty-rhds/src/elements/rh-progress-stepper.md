---
title: "Progress Steps"
imports:
  - rh-progress-stepper
---

<p>9 demos for <code>&lt;rh-progress-stepper&gt;</code></p>


### color context


<rh-context-demo>
  <rh-context-picker id="picker" target="navigation" allow="darkest, lightest"></rh-context-picker>
  <rh-progress-stepper>
    <rh-progress-step
        state="complete"
        description="This step has been completed successfully">
      Complete Step
    </rh-progress-step>
    <rh-progress-step
        state="warn"
        description="This step has a warning that needs attention">
      Warning Step
    </rh-progress-step>
    <rh-progress-step
        state="fail"
        description="This step has failed and needs to be retried">
      Failed Step
    </rh-progress-step>
    <rh-progress-step
        state="active"
        description="Currently working on this step">
      Active Step
    </rh-progress-step>
    <rh-progress-step
        state="inactive"
        description="This step is not yet started">
      Inactive Step
    </rh-progress-step>
  </rh-progress-stepper>
</rh-context-demo>



### compact horizontal


<rh-progress-stepper compact>
  <rh-progress-step
      state="complete"
      description="This step has been completed successfully">
    Complete Step
  </rh-progress-step>
  <rh-progress-step
      state="warn"
      description="This step has a warning that needs attention">
    Warning Step
  </rh-progress-step>
  <rh-progress-step
      state="fail"
      description="This step has failed and needs to be retried">
    Failed Step
  </rh-progress-step>
  <rh-progress-step
      state="active"
      description="Currently working on this step">Active Step</rh-progress-step>
  <rh-progress-step
      state="inactive"
      description="This step is not yet started">
    Inactive Step
  </rh-progress-step>
</rh-progress-stepper>



### compact vertical


<rh-progress-stepper compact orientation="vertical">
  <rh-progress-step
      state="complete"
      description="This step has been completed successfully">
    Complete Step
  </rh-progress-step>
  <rh-progress-step
      state="warn"
      description="This step has a warning that needs attention">
    Warning Step
  </rh-progress-step>
  <rh-progress-step
      state="fail"
      description="This step has failed and needs to be retried">
    Failed Step
  </rh-progress-step>
  <rh-progress-step
      state="active"
      description="Currently working on this step">
    Active Step
  </rh-progress-step>
  <rh-progress-step
      state="inactive"
      description="This step is not yet started">
    Inactive Step
  </rh-progress-step>
</rh-progress-stepper>



### custom icon


<rh-progress-stepper>
  <rh-progress-step
      state="complete"
      description="This step has been completed successfully">
    Complete Step
  </rh-progress-step>
  <rh-progress-step
      state="complete"
      description="This step has been completed successfully">
    Complete Step
  </rh-progress-step>
  <rh-progress-step
      state="complete"
      description="This step has been completed successfully">
    Current Step
  </rh-progress-step>
  <rh-progress-step
      description="Currently working on this step"
      icon="hourglass">
    Custom
  </rh-progress-step>
  <rh-progress-step
      description="This step is not yet started">
    Inactive Step
  </rh-progress-step>
</rh-progress-stepper>



### index


<rh-progress-stepper>
  <rh-progress-step state="complete"
                    description="This step has been completed successfully">
    Complete Step
  </rh-progress-step>
  <rh-progress-step state="warn"
                    description="This step has a warning that needs attention">
    Warning Step
  </rh-progress-step>
  <rh-progress-step state="fail"
                    description="This step has failed and needs to be retried">
    Failed Step
  </rh-progress-step>
  <rh-progress-step state="active"
                    description="Currently working on this step">
    Active Step
  </rh-progress-step>
  <rh-progress-step description="This step is not yet started">
    Inactive Step
  </rh-progress-step>
</rh-progress-stepper>



### linked steps


<rh-progress-stepper>
  <rh-progress-step state="complete"
                    href="#"
                    description="This step has been completed successfully">
    Complete Step
  </rh-progress-step>
  <rh-progress-step state="warn"
                    href="#"
                    description="This step has a warning that needs attention">
    Warning Step
  </rh-progress-step>
  <rh-progress-step state="fail"
                    href="#"
                    description="This step has failed and needs to be retried">
    Failed Step
  </rh-progress-step>
  <rh-progress-step state="active"
                    href="#"
                    description="Currently working on this step">
    Active Step
  </rh-progress-step>
  <rh-progress-step href="#"
                    description="This step is not yet started">
    Inactive Step
  </rh-progress-step>
</rh-progress-stepper>



### rich descriptions


<rh-progress-stepper>
  <rh-progress-step state="complete">
    Complete Step
    <span slot="description">
      This step has been <strong>completed</strong> successfully
    </span>
  </rh-progress-step>
  <rh-progress-step state="warn">
    Warning Step
    <span slot="description">
      This step has a <em>warning</em> that needs attention
    </span>
  </rh-progress-step>
  <rh-progress-step state="fail">
    Failed Step
    <span slot="description">
      This step has <strong>failed</strong> and needs to be retried
    </span>
  </rh-progress-step>
  <rh-progress-step state="active">
    Active Step
    <span slot="description">
      Currently <u>working</u> on this step
    </span>
  </rh-progress-step>
  <rh-progress-step>
    Inactive Step
    <span slot="description">
      This step is <s>not yet started</s>
    </span>
  </rh-progress-step>
</rh-progress-stepper>



### vertical at


<rh-progress-stepper vertical-at="sm">
  <rh-progress-step state="complete"
                    description="This step has been completed successfully">
    Complete Step
  </rh-progress-step>
  <rh-progress-step state="warn"
                    description="This step has a warning that needs attention">
    Warning Step
  </rh-progress-step>
  <rh-progress-step state="fail"
                    description="This step has failed and needs to be retried">
    Failed Step
  </rh-progress-step>
  <rh-progress-step state="active"
                    description="Currently working on this step">
    Active Step
  </rh-progress-step>
  <rh-progress-step description="This step is not yet started">
    Inactive Step
  </rh-progress-step>
</rh-progress-stepper>



### vertical stepper


<rh-progress-stepper orientation="vertical">
  <rh-progress-step state="complete"
                    description="This step has been completed successfully">
    Complete Step
  </rh-progress-step>
  <rh-progress-step state="warn"
                    description="This step has a warning that needs attention">
    Warning Step
  </rh-progress-step>
  <rh-progress-step state="fail"
                    description="This step has failed and needs to be retried">
    Failed Step
  </rh-progress-step>
  <rh-progress-step state="active"
                    description="Currently working on this step">
    Active Step
  </rh-progress-step>
  <rh-progress-step description="This step is not yet started">
    Inactive Step
  </rh-progress-step>
</rh-progress-stepper>


