% Author: CatInCosmicSpace
\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{diploma}[2020/06/01 Report template class]

\LoadClass{article}

\def\abstractLine{Выпускная квалификационная работа бакалавра состоит из}
\def\termsAndDefinitionsLine{В настоящей выпускной квалификационной работе бакалавра
применяют следующие термины с соответствующими определениями:}
\def\abbreviationsLine{В настоящей выпускной квалификационной работе бакалавра
применяют следующие сокращения и обозначения:}
\def\termsAndAbbreviationsLine{В настоящей выпускной квалификационной работе бакалавра
применяют следующие определения, обозначения и сокращения:}

\RequirePackage[russian]{babel}
\usepackage{fancyvrb} % в преамбулу

\RequirePackage{styles/diploma/styles/01-base}
\RequirePackage{styles/diploma/styles/02-construction}
\RequirePackage{styles/diploma/styles/03-numbering}
\RequirePackage{styles/diploma/styles/04-section-numbering}
\RequirePackage{styles/diploma/styles/05-figures}
\RequirePackage{styles/diploma/styles/06-tables}
\RequirePackage{styles/diploma/styles/07-footnotes}
\RequirePackage{styles/diploma/styles/08-formulas}
\RequirePackage{styles/diploma/styles/09-cites}
\RequirePackage{styles/diploma/styles/10-titlepage}
\RequirePackage{styles/diploma/styles/11-performers}
\RequirePackage{styles/diploma/styles/12-abstract}
\RequirePackage{styles/diploma/styles/13-contents}
\RequirePackage{styles/diploma/styles/14-terms-and-definitions}
\RequirePackage{styles/diploma/styles/15-list-of-abbreviations}
\RequirePackage{styles/diploma/styles/16-references}
\RequirePackage{styles/diploma/styles/17-appendices}
\RequirePackage{styles/diploma/styles/18-extra}
\RequirePackage{styles/diploma/styles/19-counters}
\RequirePackage{styles/diploma/styles/20-listing}
\RequirePackage{styles/diploma/styles/21-math}
\RequirePackage{styles/diploma/styles/22-algorithms}

\makeglossaries
