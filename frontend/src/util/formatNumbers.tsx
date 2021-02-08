function formatAmountTwoDecimals(amount: string): string {
  return parseFloat(amount).toFixed(2);
}

export default formatAmountTwoDecimals;
