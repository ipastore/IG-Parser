# IG-Parser
Parser for IG 2.0 Statements based on the IG Script Notation. 

Contact: Christopher Frantz (christopher.frantz@ntnu.no)

## Overview

IG-Parser is a parser for IG Script, a formal notation for the representation institutional statements (e.g., policy statements) used in the Institutional Grammar 2.0. The parser can be used by direct invocation, as well as via a web interface that produces tabular output of parsed statements (currently supporting Google Sheets format). In the following, you will find a brief introduction to syntactic principles and essential features of IG Script, followed by a set of examples showcasing all features. As a final aspect, deployment instructions for IG Parser are provided. 

The conceptual background of the Institutional Grammar 2.0 is provided [here](https://doi.org/10.1111/padm.12719), augmented with supplementary operational [coding guidelines](https://arxiv.org/abs/2008.08937).

## IG Script

IG Script is a notation introduced in the context of the [Institutional Grammar 2.0](https://doi.org/10.1111/padm.12719) (IG 2.0) that aims at a deep structural representation of legal statements alongside selected levels of expressiveness. While IG 2.0 highlights the conceptual background, the objective of IG Script is to provide an accessible, but formal approach to provide a format-independent representation of institutional statements. While the parser currently only supports export in tabular format, future refinements will include other formats (e.g.,  XML, JSON, YAML). While this introduction focuses on the operational use, syntactic and semantic foundations are provided [elsewhere](https://github.com/InstitutionalGrammar/IG-2.0-Resources).

### Principles of IG Script Syntax

IG Script centers around a set of fundamental primitives that can be combined to parse statements comprehensively, including:

* Component annotations
* Component combinations
* Nested statements
* Nested statement combinations
* Semantic Annotations (not yet supported in output)

#### Component Annotations

*Component Annotations* provide the basic building block for any statement annotation. A component is represented as 

* `componentSymbol(naturalLanguageText)`
  
`componentSymbol` is one of the supported symbols for the different component types (see below). 
An example for an annotated Attributes component is `A(Farmer)`.

`naturalLanguageText` is the human-readable text annotated as the entity the annotation describes. The text is
open-ended and can include special symbols, including parentheses (e.g., `A(Farmer (e.g., organic farmer))`). 
Exceptions to this rule are discussed in the context of combinations).

The scope of a component is specified by opening and closing parentheses.

All components of a statement are annotated correspondingly, without concern for order, or repetition. 
The parser further tolerates multiple component annotations of the same kind. Multiple Attributes, 
for example (e.g., `A(Farmer) D(must) I(comply) A(Certifier)`, are effectively interpreted as a combination (i.e., 
`A((Farmer [AND] Certifier)) D(must) I(comply)`) in the parsing process.  

Any symbols outside the annotations and combinations are ignored by the parser.

The parser supports a fixed set of component type symbols that uniquely identify a given component type.

_Supported Component Type Symbols:_

* `A` - Attributes
* `A,p` - Attributes Property
* `D` - Deontic
* `I` - Aim
* `Bdir` - Direct Object
* `Bdir,p` - Direct Object Property
* `Bind` - Indirect Object
* `Bind,p` - Indirect Object Property
* `Cac` - Activation Condition
* `Cex` - Execution Constraint
* `E` - Constituted Entity
* `E,p` - Constituted Entity Property
* `M` - Modal
* `F` - Constitutive Function 
* `P` - Constituting Properties
* `P,p` - Constituting Properties Property
* `O` - Or else

#### Component Combinations

Statements often express alternatives or combinations of activities that need to be 
administered, thus reflecting logical combinations of entities, such as actors, actions, conditions, etc. 
To combine individual components logically, IG Script supports the notion of *Component 
Combinations*. As indicated above, separate components of the same type are interpreted as 
AND-combined. To explicitly specify the nature of the logical relationship (e.g., conjunction, 
inclusive/exclusive disjunction), combinations need to be explicitly specified in the following format:

`(componentAnnotation1 [logicalOperator] componentAnnotation2)`

Example: 
* `(A(actor) [XOR] A(owner))`
* `(I(report) [XOR] I(review))`

Component combinations can be nested arbitrarily deep, i.e., any component can be a combination itself, for example:

`(componentAnnotation1 [logicalOperator] (componentAnnotation2 [logicalOperator] componentAnnotation3))`

Example:
* `(A(certifier) [AND] (A(owner) [XOR] A(inspector)))` 
* `((A(operator) [OR] A(certifier)) [AND] (A(owner) [XOR] A(inspector)))`

Supported logical operators:

* `[AND]` - conjunction (i.e., "and")
* `[OR]` - inclusive disjunction (i.e., "and/or")
* `[XOR]` - exclusive disjunction (i.e., "either or")

Invalid operators (e.g., `[AN]`) will be ignored in the parsing process.

#### Nested Statements

Selected components can be substituted by statements entirely. 
For example, the activation condition could consist of a statement 
on its own. The syntax is as follows:

`componentSymbol{ componentSymbol(naturalLanguageText) ... }`

As before, components are annotated using parentheses, and are now augmented 
with braces that delineate the nested statements.

Essentially, a fully annotated statement (e.g., `A(), I(), Cex()`) is 
framed by the statement annotation (e.g., `Cac{ A(), I(), Cex() }`).

Nesting can occur to arbitrary depth, i.e., a nested statement 
can contain another nested statement, e.g., `Cac{ A(), I(), Cac{ A(), I(), Cac() } }`, 
and can further include combinations as introduced in the following.

It is important to note that component-level nesting is limited to specific 
component types and properties.

Components for which nested statements are supported:

* `A,p` - Attributes Property
* `Bdir` - Direct Object
* `Bdir,p` - Direct Object Property
* `Bind` - Indirect Object
* `Bind,p` - Indirect Object Property
* `Cac` - Activation Condition
* `Cex` - Execution Constraint
* `E,p` - Constituted Entity Property
* `P` - Constituting Properties
* `P,p` - Constituting Properties Property
* `O` - Or else

#### Nested Statement Combinations

Nested statements can, analogous to components, be combined to an 
arbitrary depth and using the same logical operators as for component
combinations. Two constraints are to be considered: 

* Combinations of statements need to be surrounded by braces (i.e., `{` and `}`).
* Only components of the same kind can be combined.

The following example correctly reflects the combination of two nested AND-combined 
activation conditions:

`{ Cac{ A(), I(), Cex() } [AND] Cac{ A(), I(), Cex() } }`

This example, in contrast, will fail (note the differing component types `Cac` and `Cex`): 

`{ Cac{ A(), I(), Cex() } [AND] Cex{ A(), I(), Cex() } }`

Another important aspect are the outer braces surrounding the nested statement 
combination, i.e., form a `{ ... [AND] ... }` pattern (where logical operators can vary, of course).

Unlike the first example, the following one will result in an error due to missing outer braces:

`Cac{ A(), I(), Cex() } [AND] Cac{ A(), I(), Cex() }`

Note that non-nested and nested components can be used in the same statement 
(e.g., `... Cac( text ), Cac{ A(text) I(text) Cac(text) } ...` ), and are implicitly AND-combined.

#### Semantic Annotations

In addition to the parsing of component annotations and combinations of various kinds, the parser further supports semantic annotations of components according to the taxonomies outlined in the [Institutional Grammar 2.0 Codebook](https://arxiv.org/abs/2008.08937).

The syntax is `compononentSymbol[semanticAnnotation](component content)`, i.e., any component can be augmented with `[semantic annotation content]`, e.g., `Cac[context=state](Upon certification)`.

Note that semantic annotations are not yet supported by the parser.

### Examples

In the following, you will find selected examples that highlight the practical use of the features introduced above.

* Simple regulative statement: 
  * `A(Operator) D(must) I(comply) Bdir(with regulations)`
* Component combination: 
  * `A(Operator) D(must) I((comply [AND] respond)) Bdir(with/to regulations)`
* Component-level nesting: 
  * `A(Farmer) D(must) I(comply) with Bdir(provisions) Cac{A(Farmer) I(has lodged) for Bdir(application) Bdir,p(certification)}`
* Component-level nesting over multiple levels: 
  * `A(Farmer) D(must) I(comply) with Bdir(Organic Farming provisions) Cac{A(Farmer) I(has lodged) for Bdir(application) Bdir,p(certification) Cex(successfully) with the Bind(Organic Farming Program) Cac{E(Organic Program) F(covers) P,p(relevant) P(region)}}`
* Component-level nesting on various components (e.g., Bdir and Cac) and combinations of nested statements (Cac): 
  * `A(Program Manager) D(may) I(administer) Bdir(sanctions) {Cac{A(Program Manager) I(suspects) Bdir{A(farmer) I((violates [OR] does not comply)) with Bdir(regulations)}} [OR] Cac{A(Program Manager) I(has witnessed) Bdir,p(farmer's) Bdir(non-compliance) Cex(in the past)}}`
* Complex statement showcasing combined use of all features across various levels: 
  * `A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), D(may) I(inspect), I((review [AND] (revise [AND] resubmit))) Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) Cex(for compliance with the (Act or [XOR] regulations in this part)) if {Cac{A(Programme Manager) I((suspects [OR] establishes)) Bdir(violations)} [AND] Cac{E(Program Manager) F(is authorized) for the P,p(relevant) P(region)}}`
<!--* Semantic annotations: `A,p(Certified) A(organic farmers) D(must) I(submit) Bdir(report) to Bind(Organic Program Representative) Cac[context=tim](at the end of each year).`-->

* Common issues:
  * Parentheses/braces need to match. The parser calls out if a mismatch exists. 
  * Explicit specification of combinations using parentheses/braces is necessary. 
    * This example works: `A(( actor1 [XOR] actor2 ))` 
    * This one does not: `A( actor1 [XOR] actor2 )`, due to unscoped combinations.
    * The same applies for braces used for statement-level combinations, i.e., 
  `{Cac{A(actor1) I(complies) Cac(at all time)} [OR] Cac{A(actor1) I(complies) Cac(at all time)}}` will parse successfully.
  * In nested statements, only components of the same kind can be logically linked (e.g., `{Cac{ ... } [AND] Cac{ ... }}` will parse; `{Cac{ ... } [AND] Bdir{ ... }}` will not parse).

## Deployment

* Prerequisites:
  * [Docker](https://docs.docker.com/engine/install/)
  * [docker-compose](https://docs.docker.com/compose/install/) (optional if environment, volume and port parameterization is done manually)
  * Quick installation of docker under Ubuntu LTS: `sudo apt install docker.io`

* Deployment Guidelines
  * Clone this repository
  * Make `deploy.sh` executable (`chmod 740 deploy.sh`)
  * Run `deploy.sh` with superuser permissions (`sudo deploy.sh`)
    * This script deletes old versions of IG-Parser, before pulling the latest version and deploying it.
    * For manual start, run `sudo docker-compose up -d`. Run `sudo docker-compose down` to stop execution.
  
* Service Configuration
  * By default, the web service listens on port 4040, and logging is enabled in the subfolder `./logs`. 
  * The service automatically restarts if it crashes. 
  * Adjust the docker-compose.yml file to modify any of these characteristics.
