import * as React from "react";

interface BankModalProps {
  showBankModal: boolean;
  setShowBankModal: (showBankModal: boolean) => void;
}

function BankModal(props: BankModalProps): JSX.Element | null {
  if (!props.showBankModal) {
    return null;
  }

  return (
    <div
      className="modal fade"
      id="chooseBankAccountModal"
      tabIndex={-1}
      role="dialog"
      aria-labelledby="chooseBankAccountModalTitle"
      aria-hidden="true"
    >
      <div className="modal-dialog modal-dialog-centered" role="document">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="modal-title" id="chooseBankAccountModalTitle">
              Choose bank account
            </h5>
            <button
              type="button"
              className="close"
              data-dismiss="modal"
              aria-label="Close"
              onClick={() => props.setShowBankModal(false)}
            >
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div className="modal-body">
            <div className="container-fluid">
              <div className="row">
                <div className="col-md-6 d-flex justify-content-center">
                  <button
                    type="button"
                    className="btn btn-primary"
                    id="getTransactions"
                    data-dismiss="modal"
                  >
                    <img
                      width="15px"
                      height="15px"
                      src="https://monzo.com/static/images/mondo-mark-01.png"
                    ></img>
                    Monzo
                  </button>
                </div>
                <div className="col-md-6 d-flex justify-content-center">
                  <button
                    type="button"
                    className="btn btn-primary mx-auto"
                    data-dismiss="modal"
                  >
                    <img
                      width="15px"
                      height="15px"
                      src="https://alternative.me/media/256/revolut-icon-3t64wiq24kxp057j-c.png"
                    ></img>
                    Revolut
                  </button>
                </div>
              </div>
              <div className="row">
                <div className="col-md-6 d-flex justify-content-center">
                  <button
                    type="button"
                    className="btn btn-primary"
                    data-dismiss="modal"
                  >
                    <img
                      width="15px"
                      height="15px"
                      src="https://play-lh.googleusercontent.com/1U3nHP3cS5s8yNuIH4ECo-5bi_lUJ4dZyxO2HPCZSrlPeVAE5UQSIszDt__3fv36GK8"
                    ></img>
                    Barclays
                  </button>
                </div>
                <div className="col-md-6 d-flex justify-content-center">
                  <button
                    type="button"
                    className="btn btn-primary"
                    data-dismiss="modal"
                  >
                    <img
                      width="15px"
                      height="15px"
                      src="https://appmirror.net/wp-content/uploads/2020/11/natwest-icon-1200x1200.png"
                    ></img>
                    Natwest
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default BankModal;
