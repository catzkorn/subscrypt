import Subscription from "../subscriptions/subscriptionType";

function sortDates(a: Subscription, b: Subscription): number {
  const dayA = new Date(a.dateDue).getDate();
  const dayB = new Date(b.dateDue).getDate();
  if (dayA < dayB) {
    return -1;
  } else if (dayA > dayB) {
    return 1;
  }
  return 0;
}

export default sortDates;
