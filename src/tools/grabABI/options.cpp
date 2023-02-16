/*-------------------------------------------------------------------------------------------
 * qblocks - fast, easily-accessible, fully-decentralized data from blockchains
 * copyright (c) 2016, 2021 TrueBlocks, LLC (http://trueblocks.io)
 *
 * This program is free software: you may redistribute it and/or modify it under the terms
 * of the GNU General Public License as published by the Free Software Foundation, either
 * version 3 of the License, or (at your option) any later version. This program is
 * distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
 * the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details. You should have received a copy of the GNU General
 * Public License along with this program. If not, see http://www.gnu.org/licenses/.
 *-------------------------------------------------------------------------------------------*/
/*
 * Parts of this file were generated with makeClass --options. Edit only those parts of
 * the code outside of the BEG_CODE/END_CODE sections
 */
#include "options.h"

//---------------------------------------------------------------------------------------------------
static const COption params[] = {
    // BEG_CODE_OPTIONS
    // clang-format off
    COption("addrs", "", "list<addr>", OPT_REQUIRED | OPT_POSITIONAL, "a list of one or more smart contracts whose ABIs to display"),  // NOLINT
    COption("known", "k", "", OPT_SWITCH, "load common 'known' ABIs from cache"),
    COption("", "", "", OPT_DESCRIPTION, "Fetches the ABI for a smart contract."),
    // clang-format on
    // END_CODE_OPTIONS
};
static const size_t nParams = sizeof(params) / sizeof(COption);

//---------------------------------------------------------------------------------------------------
bool COptions::parseArguments(string_q& command) {
    if (!standardOptions(command))
        return false;

    // BEG_CODE_LOCAL_INIT
    bool known = false;
    // END_CODE_LOCAL_INIT

    Init();
    explode(arguments, command, ' ');
    blknum_t latest = NOPOS;  // getLatestBlock_client();
    string_q unused;
    for (auto arg : arguments) {
        if (false) {
            // do nothing -- make auto code generation easier
        } else if (isKnownAbi(arg, unused)) {
            addrs.push_back(arg);

            // BEG_CODE_AUTO
        } else if (arg == "-k" || arg == "--known") {
            known = true;

        } else if (arg == "-s" || arg == "--sol") {
            // clang-format off
            return usage("the --sol option is deprecated, please use the `solc --abi` tool instead");  // NOLINT
            // clang-format on

        } else if (startsWith(arg, '-')) {  // do not collapse

            if (!builtInCmd(arg)) {
                return invalid_option(arg);
            }

        } else {
            if (!parseAddressList(this, addrs, arg))
                return false;

            // END_CODE_AUTO
        }
    }

    if (known)
        abi_spec.loadAbisFromKnown();

    for (auto addr : addrs) {
        bool testing = isTestMode() && addr == "0xeeeeeeeeddddddddeeeeeeeeddddddddeeeeeeee";
        string_q fileName = cacheFolder_abis + addr + ".json";
        if (!testing && !isContractAt(addr, latest) && !fileExists(fileName)) {
            cerr << "Address " << addr << " is not a smart contract. Skipping..." << endl;
        } else {
            abi_spec.loadAbiFromEtherscan(addr);
            abi_spec.address = addr;
        }
    }

    // Display formatting
    string_q format = STR_DISPLAY_ABI;
    string_q ffields = toLower(substitute(substitute(substitute(STR_DISPLAY_FUNCTION, "[{", ""), "}]", ""), "\t", ","));
    string_q funcFields = "CFunction:" + ffields;
    if (isTestMode())
        funcFields = "CFunction:" + substitute(ffields, "inputs,outputs", "input_names,output_names");
    replace(format, "[{ADDRESS}]\t", "");
    replace(funcFields, "address,", "");
    replace(format, "[{ABI_SOURCE}]\t", "");
    replace(funcFields, "abi_source,", "");

    if (verbose && (expContext().exportFmt == JSON1 || isApiMode())) {
        replaceAll(funcFields, "_name", "");
        replaceAll(format, "_NAME", "");
    }

    configureDisplay("grabABI", "CAbi", format);
    manageFields("CFunction:all", false);
    manageFields(funcFields, true);
    manageFields("CParameter:all", false);
    manageFields("CParameter:type,name,internalType,components,is_array,indexed", true);

    return true;
}

//---------------------------------------------------------------------------------------------------
void COptions::Init(void) {
    // BEG_CODE_GLOBALOPTS
    registerOptions(nParams, params, OPT_RAW, OPT_DENOM);
    // END_CODE_GLOBALOPTS

    // BEG_CODE_INIT
    // END_CODE_INIT

    parts = (SIG_DEFAULT | SIG_ENCODE);
    addrs.clear();
}

//---------------------------------------------------------------------------------------------------
COptions::COptions(void) {
    Init();

    // BEG_CODE_NOTES
    // clang-format off
    notes.push_back("Search for either four byte signatures or event signatures with the --find option.");
    // clang-format on
    // END_CODE_NOTES

    // BEG_ERROR_STRINGS
    // END_ERROR_STRINGS
}

//--------------------------------------------------------------------------------
COptions::~COptions(void) {
}
