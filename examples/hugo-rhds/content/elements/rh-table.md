---
title: "Table"
imports:
  - rh-table
lightdom:
  - rh-table-lightdom.css
---

<p>12 demos for <code>&lt;rh-table&gt;</code></p>


### auto data labels

<h2>Basic table</h2>
<rh-table>
  <table>
    <col>
    <col>
    <col>
    <thead>
      <tr>
        <th scope="col">Date</th>
        <th scope="col">Event</th>
        <th scope="col">Venue</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>12 February</td>
        <td>
          <a href="">Waltz with Strauss</a>
        </td>
        <td>Main Hall</td>
      </tr>
      <tr>
        <td>24 March</td>
        <td>The Obelisks</td>
        <td>West Wing</td>
      </tr>
      <tr>
        <td data-label="Custom 1">14 April</td>
        <td data-label="Custom 2">The What</td>
        <td data-label="Custom 3">Main Hall</td>
      </tr>
    </tbody>
  </table>
</rh-table>
<h2>Table with <code>rowspans</code></h2>
<rh-table>
  <table>
    <col>
    <col>
    <col>
    <thead>
      <tr>
        <th scope="col">Date</th>
        <th scope="col">Event</th>
        <th scope="col">Venue</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>12 February</td>
        <td rowspan="2">
          <a href="">Waltz with Strauss</a>
        </td>
        <td>Main Hall</td>
      </tr>
      <tr>
        <td>24 March</td>
        <td>West Wing</td>
      </tr>
      <tr>
        <td data-label="Custom 1">14 April</td>
        <td data-label="Custom 2">The What</td>
        <td data-label="Custom 3">Main Hall</td>
      </tr>
    </tbody>
  </table>
</rh-table>
<h2>Table with <code>colspans</code></h2>
<rh-table>
  <table>
    <col>
    <col>
    <col>
    <thead>
      <tr>
        <th scope="col">Date</th>
        <th scope="col">Event</th>
        <th scope="col">Venue</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>12 February</td>
        <td>
          <a href="">Waltz with Strauss</a>
        </td>
        <td>Main Hall</td>
      </tr>
      <tr>
        <td>24 March</td>
        <td colspan="2">The Obelisks in the West Wing</td>
      </tr>
      <tr>
        <td>14 April</td>
        <td>The What</td>
        <td>Main Hall</td>
      </tr>
    </tbody>
  </table>
</rh-table>
<h2>Table with a complicated <code>thead</code></h2>
<rh-table>
  <table>
    <col>
    <col>
    <col>
    <thead>
      <tr>
        <th scope="col"></th>
        <th id="th-when-where" scope="col" colspan="2">When &amp; where</th>
      </tr>
      <tr>
        <th id="th-event" scope="col">Event</th>
        <th id="th-date" scope="col">Date</th>
        <th id="th-venue" scope="col">Venue</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <th headers="th-event" id="th-event-waltz" scope="row">
          <a href="">Waltz with Strauss</a>
        </td>
        <td headers="th-event-waltz th-when-where th-date">12 February</td>
        <td headers="th-event-waltz th-when-where th-venue">Main Hall</td>
      </tr>
      <tr>
        <th headers="th-event" id="th-event-obelisks" scope="row">The Obelisks</td>
        <td headers="th-event-obelisks th-when-where th-date">24 March</td>
        <td headers="th-event-obelisks th-when-where th-venue">West Wing</td>
      </tr>
      <tr>
        <th headers="th-event" id="th-event-thewhat" scope="row" data-label="Custom event">The What</td>
        <td headers="th-event-thewhat th-when-where th-date" data-label="Custom date">14 April</td>
        <td headers="th-event-thewhat th-when-where th-venue" data-label="Custom location">Main Hall</td>
      </tr>
    </tbody>
  </table>
</rh-table>


### color context

<rh-context-demo>
  <rh-table>
    <table>
      <caption>Concerts</caption>
      <colgroup>
        <col>
        <col>
        <col>
      </colgroup>
      <thead>
        <tr>
          <th scope="col">Date</th>
          <th scope="col">Event<rh-sort-button></rh-sort-button></th>
          <th scope="col">Venue<rh-sort-button></rh-sort-button></th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>12 February</td>
          <td>Waltz with Strauss</td>
          <td>Main Hall</td>
        </tr>
        <tr>
          <td>24 March</td>
          <td>The Obelisks</td>
          <td>West Wing</td>
        </tr>
        <tr>
          <td>14 April</td>
          <td>The What</td>
          <td><a href="#">Main</a> Hall</td>
        </tr>
      </tbody>
    </table>
    <small slot="summary">Dates and venues subject to change.</small>
  </rh-table>
</rh-context-demo>


### column headers

<p>
  <strong>Note: </strong> Tables with no <code>thead</code> will not stack on mobile.
</p>
<rh-table>
  <table>
    <tbody>
      <tr>
        <th scope="row">Date</th>
        <td>12 February</td>
        <td>24 March</td>
        <td>14 April</td>
      </tr>
      <tr>
        <th scope="row">Event</th>
        <td>Waltz with Strauss</td>
        <td>The Obelisks</td>
        <td>The What</td>
      </tr>
      <tr>
        <th scope="row">Venue</th>
        <td>Main Hall</td>
        <td>West Wing</td>
        <td>Main Hall</td>
      </tr>
    </tbody>
  </table>
</rh-table>


### headers and summary but no title

<rh-table>
  <table>
    <col>
    <col>
    <col>
    <thead>
      <tr>
        <th scope="col" data-label="Date">Date</th>
        <th scope="col" data-label="Event">Event</th>
        <th scope="col" data-label="Venue">Venue</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td data-label="Date">12 February</td>
        <td data-label="Event">
          <a href="">Waltz with Strauss</a>
        </td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Date">24 March</td>
        <td data-label="Event">The Obelisks</td>
        <td data-label="Venue">West Wing</td>
      </tr>
      <tr>
        <td data-label="Date">14 April</td>
        <td data-label="Event">The What</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
    </tbody>
  </table>
</rh-table>


### horizontal overflow

<rh-table style="width:500px">
  <table>
    <tbody>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
    </tbody>
  </table>
</rh-table>


### index

<rh-table>
  <table>
    <caption>
      Concerts
    </caption>
    <colgroup>
      <col>
      <col>
      <col>
    </colgroup>
    <thead>
      <tr>
        <th scope="col">Date</th>
        <th scope="col">Event<rh-sort-button></rh-sort-button>
        </th>
        <th scope="col">Venue<rh-sort-button></rh-sort-button>
        </th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>12 February</td>
        <td>Waltz with Strauss</td>
        <td>Main Hall</td>
      </tr>
      <tr>
        <td>24 March</td>
        <td>The Obelisks</td>
        <td>West Wing</td>
      </tr>
      <tr>
        <td>14 April</td>
        <td>The What</td>
        <td>Main Hall</td>
      </tr>
    </tbody>
  </table>
  <small slot="summary">Dates and venues subject to change.</small>
</rh-table>


### no title headers or summary

<rh-table>
  <table>
    <col>
    <col>
    <col>
    <tbody>
      <tr>
        <td data-label="Date">12 February</td>
        <td data-label="Event">
          <a href="">Waltz with Strauss</a>
        </td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Date">24 March</td>
        <td data-label="Event">The Obelisks</td>
        <td data-label="Venue">West Wing</td>
      </tr>
      <tr>
        <td data-label="Date">14 April</td>
        <td data-label="Event">The What</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
    </tbody>
  </table>
</rh-table>


### row and column headers

<rh-table>
  <table>
    <caption>Delivery slots:</caption>
    <col>
    <col>
    <col>
    <col>
    <col>
    <col>
    <thead>
      <tr>
        <td></td>
        <th sortable scope="col" data-label="Monday">Monday</th>
        <th sortable scope="col" data-label="Tuesday">Tuesday</th>
        <th sortable scope="col" data-label="Wednesday">Wednesday</th>
        <th sortable scope="col" data-label="Thursday">Thursday</th>
        <th sortable scope="col" data-label="Friday">Friday</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <th scope="row" data-label="Time">09:00 – 11:00</th>
        <td data-label="Monday">Closed</td>
        <td data-label="Tuesday">Open</td>
        <td data-label="Wednesday">Open</td>
        <td data-label="Thursday">Closed</td>
        <td data-label="Friday">Closed</td>
      </tr>
      <tr>
        <th scope="row" data-label="Time">11:00 – 13:00</th>
        <td data-label="Monday">Open</td>
        <td data-label="Tuesday">Open</td>
        <td data-label="Wednesday">Closed</td>
        <td data-label="Thursday">Closed</td>
        <td data-label="Friday">Closed</td>
      </tr>
      <tr>
        <th cope="row" data-label="Time">13:00 – 15:00</th>
        <td data-label="Monday">Open</td>
        <td data-label="Tuesday">Open</td>
        <td data-label="Wednesday">Open</td>
        <td data-label="Thursday">Closed</td>
        <td data-label="Friday">Closed</td>
      </tr>
      <tr>
        <th scope="row" data-label="Time">15:00 – 17:00</th>
        <td data-label="Monday">Closed</td>
        <td data-label="Tuesday">Closed</td>
        <td data-label="Wednesday">Closed</td>
        <td data-label="Thursday">Open</td>
        <td data-label="Friday">Open</td>
      </tr>
    </tbody>
  </table>
</rh-table>


### row headers

<rh-table>
  <table>
    <col>
    <col>
    <col>
    <thead>
      <tr>
        <th sortable scope="col" data-label="Date">Date</th>
        <th sortable scope="col" data-label="Event">Event</th>
        <th sortable scope="col" data-label="Venue">Venue</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td data-label="Date">12 February</td>
        <td data-label="Event">
          <a href="">Waltz with Strauss</a>
        </td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Date">24 March</td>
        <td data-label="Event">The Obelisks</td>
        <td data-label="Venue">West Wing</td>
      </tr>
      <tr>
        <td data-label="Date">14 April</td>
        <td data-label="Event">The What</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
    </tbody>
  </table>
</rh-table>


### title and summary but no headers

<rh-table>
  <table>
    <caption>Concerts</caption>
    <col>
    <col>
    <col>
    <tbody>
      <tr>
        <td data-label="Date">12 February</td>
        <td data-label="Event">
          <a href="">Waltz with Strauss</a>
        </td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Date">24 March</td>
        <td data-label="Event">The Obelisks</td>
        <td data-label="Venue">West Wing</td>
      </tr>
      <tr>
        <td data-label="Date">14 April</td>
        <td data-label="Event">The What</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
    </tbody>
  </table>
  <small slot="summary">Dates and venues subject to change. Anything longer should wrap.</small>
</rh-table>


### title headers and summary

<rh-table>
  <table>
    <caption>Concerts</caption>
    <colgroup>
      <col>
      <col>
      <col>
    </colgroup>
    <thead>
      <tr>
        <th scope="col" data-label="Date">Date</th>
        <th scope="col" data-label="Event">Event<rh-sort-button></rh-sort-button>
        </th>
        <th scope="col" data-label="Venue">Venue<rh-sort-button></rh-sort-button>
        </th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td data-label="Date">12 February</td>
        <td data-label="Event">
          <a href="#">Waltz with Strauss</a>
        </td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Date">24 March</td>
        <td data-label="Event">The Obelisks</td>
        <td data-label="Venue">West Wing</td>
      </tr>
      <tr>
        <td data-label="Date">14 April</td>
        <td data-label="Event">The What</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
    </tbody>
  </table>
  <small slot="summary">Dates and venues subject to change.</small>
</rh-table>


### vertical overflow

<rh-table style="height:500px">
  <table>
    <tbody>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
      <tr>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
        <td data-label="Venue">Main Hall</td>
      </tr>
    </tbody>
  </table>
</rh-table>

